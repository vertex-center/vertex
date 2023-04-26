package services

import (
	"context"
	"errors"
	"fmt"
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
	releaseVersion := *release.TagName
	releaseVersion = strings.TrimPrefix(releaseVersion, "v")

	if d.currentVersion == releaseVersion {
		logger.Log("vertex is already up-to-date").Print()
		return nil, nil
	}

	d.update = &types.Update{
		ID:             d.GetID(),
		Name:           "Vertex",
		CurrentVersion: d.currentVersion,
		LatestVersion:  releaseVersion,
	}

	logger.Log("a new release for Vertex is available").
		AddKeyValue("current", d.currentVersion).
		AddKeyValue("release", releaseVersion).
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
