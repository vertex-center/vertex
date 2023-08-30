package storage

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v50/github"
	"github.com/vertex-center/vertex/pkg/logger"
)

const (
	PathClient    = "live/client"
	PathPackages  = "live/packages"
	PathProxy     = "live/proxy"
	PathInstances = "live/instances"
	PathServices  = "live/services"
	PathSettings  = "live/settings"
	PathUpdates   = "live/updates"
)

var (
	ErrNoReleasesPublished = errors.New("this repository has no existing releases")
	ErrNoReleasesForThisOS = errors.New("this repository has no releases appropriate for this OS")
)

func CloneOrPullRepository(url string, dest string) error {
	logger.Log("downloading repository").
		AddKeyValue("url", url).
		Print()

	_, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		repo, err := git.PlainOpen(dest)
		if err != nil {
			return err
		}

		worktree, err := repo.Worktree()
		if err != nil {
			return err
		}

		err = worktree.Pull(&git.PullOptions{})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func DownloadLatestGithubRelease(owner string, repo string, dest string) error {
	logger.Log("downloading repository").
		AddKeyValue("owner", owner).
		AddKeyValue("repo", repo).
		Print()

	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		return ErrNoReleasesPublished
	}

	return DownloadGithubRelease(release, dest)
}

func DownloadGithubRelease(release *github.RepositoryRelease, dest string) error {
	logger.Log("downloading GitHub release").
		AddKeyValue("release", *release.Name).
		Print()

	platform := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)

	for _, asset := range release.Assets {
		if strings.Contains(*asset.Name, platform) {
			archivePath := path.Join(dest, "temp.tar.gz")

			err := Download(*asset.BrowserDownloadURL, dest, "temp.tar.gz")
			if err != nil {
				return err
			}

			err = UntarFile(dest, archivePath)
			if err != nil {
				return err
			}

			err = os.Remove(archivePath)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return ErrNoReleasesForThisOS
}

func Download(url string, dest string, filename string) error {
	logger.Log("downloading repository").
		AddKeyValue("url", url).
		Print()

	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return err
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(path.Join(dest, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func UntarFile(basePath string, archivePath string) error {
	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	stream, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}
	defer stream.Close()

	reader := tar.NewReader(stream)

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		filepath := path.Join(basePath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filepath, os.ModePerm)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			err := os.MkdirAll(path.Dir(filepath), os.ModePerm)
			if err != nil {
				return err
			}

			file, err := os.Create(filepath)
			if err != nil {
				return err
			}

			_, err = io.Copy(file, reader)
			if err != nil {
				return err
			}

			err = os.Chmod(filepath, 0755)
			if err != nil {
				return err
			}

			file.Close()
		default:
			return fmt.Errorf("unknown flag type (%b) for file '%s'", header.Typeflag, header.Name)
		}
	}

	return nil
}
