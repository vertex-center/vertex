package storage

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v50/github"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vertex/server/pkg/varchiver"
	"github.com/vertex-center/vertex/server/pkg/vdownloader"
	"github.com/vertex-center/vlog"
)

var ErrNoReleasesForThisOS = errors.New("this repository has no releases appropriate for this OS")

func CloneRepository(url string, dest string) error {
	log.Info("cloning repository",
		vlog.String("url", url),
	)
	_, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return err
}

func CloneOrPullRepository(url string, dest string) error {
	err := CloneRepository(url, dest)
	if err != nil && errors.Is(err, git.ErrRepositoryAlreadyExists) {
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

func DownloadGithubRelease(release *github.RepositoryRelease, dest string) error {
	log.Info("downloading release",
		vlog.String("release", *release.Name),
	)

	platform := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)

	for _, asset := range release.Assets {
		if strings.Contains(*asset.Name, platform) {
			archivePath := path.Join(dest, "temp.tar.gz")
			url := *asset.BrowserDownloadURL

			err := vdownloader.Download(url, dest, "temp.tar.gz")
			if err != nil {
				return err
			}

			err = varchiver.Untar(archivePath, dest)
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
