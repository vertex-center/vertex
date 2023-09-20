package services

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v50/github"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
	"golang.org/x/exp/slices"
)

var (
	Dependencies = []*types.Dependency{
		{
			ID:      "vertex-webui",
			Name:    "Vertex Web UI",
			Updater: &clientUpdater{},
		},
		{
			ID:      "vertex-services",
			Name:    "Vertex Services",
			Updater: newGitHubUpdater("services", "Vertex Services", "vertex-services"),
		},
		{
			ID:      "vertex-dependencies",
			Name:    "Vertex Dependencies",
			Updater: newGitHubUpdater("packages", "Vertex Dependencies", "vertex-dependencies"),
		},
	}
)

var (
	ErrDependencyNotInstalled = errors.New("dependency is not installed")
)

type DependenciesService struct {
	dependencies types.Dependencies
}

func NewDependenciesService(currentVertexVersion string) DependenciesService {
	dependencies := append(Dependencies, &types.Dependency{
		ID:      "vertex",
		Name:    "Vertex",
		Updater: &vertexUpdater{currentVersion: currentVertexVersion},
	})

	return DependenciesService{
		dependencies: types.Dependencies{
			Items: dependencies,
		},
	}
}

func (s *DependenciesService) GetCachedUpdates() types.Dependencies {
	return s.dependencies
}

func (s *DependenciesService) CheckForUpdates() (types.Dependencies, error) {
	log.Info("fetching all updates...")

	for _, dependency := range s.dependencies.Items {
		log.Info("fetching dependency",
			vlog.String("id", dependency.ID),
		)

		update, err := dependency.Updater.CheckForUpdate()
		if err != nil {
			log.Error(err)
		}
		dependency.Update = update
	}

	t := time.Now()
	s.dependencies.LastUpdatesCheck = &(t)
	return s.dependencies, nil
}

func (s *DependenciesService) InstallUpdates(dependenciesID []string) error {
	for _, dependency := range s.dependencies.Items {
		if slices.Contains(dependenciesID, dependency.ID) {
			err := dependency.Updater.InstallUpdate()
			if err != nil {
				return err
			}
			dependency.Update = nil
		}
	}
	return nil
}

// Vertex: https://github.com/vertex-center/vertex
type vertexUpdater struct {
	// The current version of Vertex.
	currentVersion string

	// The GitHub release associated with the update
	release *github.RepositoryRelease
}

func (d *vertexUpdater) CheckForUpdate() (*types.DependencyUpdate, error) {
	if d.currentVersion == "dev" {
		log.Info("skipping vertex update in 'dev' version")
		return nil, nil
	}

	// remove previous old version if it exists.
	err := os.Remove("vertex-old")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	client := github.NewClient(nil)

	// get the latest release
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "vertex-center", "vertex")
	if err != nil {
		return nil, err
	}
	d.release = release

	// check if the version is different
	latestVersion := *release.TagName
	latestVersion = strings.TrimPrefix(latestVersion, "v")

	if d.currentVersion == latestVersion {
		return nil, nil
	}

	log.Info("a new release for Vertex is available",
		vlog.String("current", d.currentVersion),
		vlog.String("release", latestVersion),
	)

	return &types.DependencyUpdate{
		CurrentVersion: d.currentVersion,
		LatestVersion:  latestVersion,
		NeedsRestart:   true,
	}, nil
}

func (d *vertexUpdater) InstallUpdate() error {
	if d.release == nil {
		return errors.New("the release has not been fetched before installing the update")
	}

	dir := path.Join(storage.Path, "updates")

	err := storage.DownloadGithubRelease(d.release, dir)
	if err != nil {
		return err
	}

	err = os.Rename("vertex", "vertex-old")
	if err != nil {
		return fmt.Errorf("failed to rename old executable: %v", err)
	}

	err = os.Rename(path.Join(dir, "vertex"), "vertex")
	if err != nil {
		return err
	}

	d.release = nil

	log.Warn("a new Vertex update has been installed. please restart Vertex to apply changes.")

	return nil
}

func (d *vertexUpdater) GetPath() string {
	return "."
}

type clientUpdater struct {
	currentVersion string
	release        *github.RepositoryRelease
}

func (d *clientUpdater) CheckForUpdate() (*types.DependencyUpdate, error) {
	d.FetchCurrentVersion()

	client := github.NewClient(nil)

	owner := "vertex-center"
	repo := "vertex-webui"

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the latest github release for %s: %v", repo, err)
	}

	latestVersion := *release.TagName

	if d.currentVersion == latestVersion {
		return nil, nil
	}

	d.release = release
	return &types.DependencyUpdate{
		CurrentVersion: d.currentVersion,
		LatestVersion:  latestVersion,
		NeedsRestart:   false,
	}, nil
}

func (d *clientUpdater) InstallUpdate() error {
	log.Info("downloading vertex-webui client...")

	for _, asset := range d.release.Assets {
		if strings.Contains(*asset.Name, "vertex-webui") {
			dir := d.GetPath()

			err := os.RemoveAll(dir)
			if err != nil {
				return err
			}

			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}

			err = download(dir, *asset.BrowserDownloadURL)
			if err != nil {
				return err
			}

			err = unarchive(dir)
			if err != nil {
				return err
			}

			err = os.Remove(path.Join(dir, "temp.zip"))
			if err != nil {
				return err
			}

			err = config.Current.Apply()
			if err != nil {
				return err
			}
		}
	}

	d.FetchCurrentVersion()
	d.release = nil

	return nil
}

func (d *clientUpdater) GetPath() string {
	return path.Join(storage.Path, "client")
}

func (d *clientUpdater) FetchCurrentVersion() {
	version, err := os.ReadFile(path.Join(d.GetPath(), "dist", "version.txt"))
	if err != nil {
		return
	}
	d.currentVersion = strings.TrimSpace(string(version))
}

func download(dir string, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(path.Join(dir, "temp.zip"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func unarchive(dir string) error {
	reader, err := zip.OpenReader(path.Join(dir, "temp.zip"))
	if err != nil {
		return err
	}

	for _, header := range reader.File {
		filepath := path.Join(dir, header.Name)

		if header.FileInfo().IsDir() {
			err = os.MkdirAll(filepath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(path.Dir(filepath), os.ModePerm)
			if err != nil {
				return err
			}

			file, err := os.Create(filepath)
			if err != nil {
				return err
			}

			content, err := header.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(file, content)
			if err != nil {
				return err
			}

			err = os.Chmod(filepath, 0755)
			if err != nil {
				return err
			}

			file.Close()
		}
	}

	return nil
}

type gitHubUpdater struct {
	dir  string
	name string
	repo string
}

func newGitHubUpdater(dir string, name string, repo string) *gitHubUpdater {
	return &gitHubUpdater{
		dir:  path.Join(storage.Path, dir),
		name: name,
		repo: repo,
	}
}

func (d *gitHubUpdater) CheckForUpdate() (*types.DependencyUpdate, error) {
	client := github.NewClient(nil)

	// Local
	repo, err := git.PlainOpen(d.dir)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		return nil, ErrDependencyNotInstalled
	}
	if err != nil {
		return nil, err
	}
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	localSHA := ref.Hash().String()

	// Remote
	branch, _, err := client.Repositories.GetBranch(context.Background(), "vertex-center", d.repo, "main", false)
	if err != nil {
		return nil, err
	}
	remoteSHA := branch.Commit.GetSHA()
	if remoteSHA == "" {
		return nil, errors.New("commit sha not found")
	}

	// Comparison
	if localSHA != remoteSHA {
		return &types.DependencyUpdate{
			CurrentVersion: localSHA,
			LatestVersion:  remoteSHA,
			NeedsRestart:   false,
		}, nil
	}

	return nil, nil
}

func (d *gitHubUpdater) InstallUpdate() error {
	url := "https://github.com/vertex-center/" + d.repo
	return storage.CloneOrPullRepository(url, d.dir)
}

func (d *gitHubUpdater) GetPath() string {
	return d.dir
}
