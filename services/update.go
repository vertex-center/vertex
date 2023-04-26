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

	"github.com/google/go-github/v50/github"
	"github.com/vertex-center/vertex/logger"
	"github.com/vertex-center/vertex/storage"
	"github.com/vertex-center/vertex/types"
	"golang.org/x/exp/slices"
)

type UpdateDependenciesService struct {
	dependencies []types.Dependency
}

func NewUpdateDependenciesService(currentVertexVersion string) UpdateDependenciesService {
	return UpdateDependenciesService{
		dependencies: []types.Dependency{
			&vertexDependency{currentVersion: currentVertexVersion},
			&VertexClientDependency{},
		},
	}
}

func (s UpdateDependenciesService) CheckForUpdates() ([]types.Update, error) {
	var updates []types.Update

	for _, dependency := range s.dependencies {
		update, err := dependency.CheckForUpdate()
		if err != nil {
			return nil, err
		}
		if update != nil {
			updates = append(updates, *update)
		}
	}

	return updates, nil
}

func (s UpdateDependenciesService) InstallUpdates(dependenciesID []string) error {
	for _, dependency := range s.dependencies {
		if slices.Contains(dependenciesID, dependency.GetID()) {
			err := dependency.InstallUpdate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Vertex: https://github.com/vertex-center/vertex
type vertexDependency struct {
	// The current version of Vertex.
	currentVersion string

	// The update found by the CheckForUpdate function.
	update *types.Update
	// The GitHub release associated with the update
	release *github.RepositoryRelease
}

func (d *vertexDependency) CheckForUpdate() (*types.Update, error) {
	if d.currentVersion == "dev" {
		logger.Log("skipping update in 'dev' version").Print()
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

	d.update = &types.Update{
		ID:             d.GetID(),
		Name:           "Vertex",
		CurrentVersion: d.currentVersion,
		LatestVersion:  latestVersion,
		NeedsRestart:   true,
	}

	logger.Log("a new release for Vertex is available").
		AddKeyValue("current", d.currentVersion).
		AddKeyValue("release", latestVersion).
		Print()

	return d.update, nil
}

func (d *vertexDependency) InstallUpdate() error {
	if d.release == nil {
		return errors.New("the release has not been fetched before installing the update")
	}

	err := storage.DownloadGithubRelease(d.release, storage.PathUpdates)
	if err != nil {
		return err
	}

	err = os.Rename("vertex", "vertex-old")
	if err != nil {
		return fmt.Errorf("failed to rename old executable: %v", err)
	}

	err = os.Rename(path.Join(storage.PathUpdates, "vertex"), "vertex")
	if err != nil {
		return err
	}

	d.currentVersion = d.update.LatestVersion
	d.release = nil
	d.update = nil

	logger.Warn("a new Vertex update has been installed. please restart Vertex to apply changes.").Print()

	return nil
}

func (d *vertexDependency) GetID() string {
	return "vertex"
}

type VertexClientDependency struct {
	currentVersion string
	release        *github.RepositoryRelease
	update         *types.Update
}

func (d *VertexClientDependency) CheckForUpdate() (*types.Update, error) {
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
	d.update = &types.Update{
		ID:             d.GetID(),
		Name:           "Vertex Client",
		CurrentVersion: d.currentVersion,
		LatestVersion:  latestVersion,
		NeedsRestart:   false,
	}
	return d.update, nil
}

func (d *VertexClientDependency) InstallUpdate() error {
	logger.Log("downloading vertex-webui client...").Print()

	for _, asset := range d.release.Assets {
		if strings.Contains(*asset.Name, "vertex-webui") {
			err := os.RemoveAll(storage.PathClient)
			if err != nil {
				return err
			}

			err = os.MkdirAll(storage.PathClient, os.ModePerm)
			if err != nil {
				return err
			}

			err = download(*asset.BrowserDownloadURL)
			if err != nil {
				return err
			}

			err = unarchive()
			if err != nil {
				return err
			}

			err = os.Remove(path.Join(storage.PathClient, "temp.zip"))
			if err != nil {
				return err
			}
		}
	}

	d.FetchCurrentVersion()
	d.update = nil
	d.release = nil

	return nil
}

func (d *VertexClientDependency) GetID() string {
	return "vertex-webui"
}

func (d *VertexClientDependency) FetchCurrentVersion() {
	version, err := os.ReadFile(path.Join(storage.PathClient, "dist", "version.txt"))
	if err != nil {
		return
	}
	d.currentVersion = strings.TrimSpace(string(version))
}

func download(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(path.Join(storage.PathClient, "temp.zip"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func unarchive() error {
	reader, err := zip.OpenReader(path.Join(storage.PathClient, "temp.zip"))
	if err != nil {
		return err
	}

	for _, header := range reader.File {
		filepath := path.Join(storage.PathClient, header.Name)

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
