package services

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/vertex-center/vertex/logger"
	"github.com/vertex-center/vertex/storage"
	"github.com/vertex-center/vertex/types"
)

type UpdateService struct{}

func NewUpdateService() UpdateService {
	return UpdateService{}
}

func (s UpdateService) CheckForUpdates(currentVertexVersion string) ([]types.Update, error) {
	var updates []types.Update

	vertexUpdate, _, err := s.CheckForVertexUpdate(currentVertexVersion)
	if err != nil {
		return nil, err
	} else {
		updates = append(updates, *vertexUpdate)
	}

	return updates, nil
}

func (s UpdateService) CheckForVertexUpdate(currentVersion string) (*types.Update, *github.RepositoryRelease, error) {
	update := &types.Update{
		Id:             "vertex",
		Name:           "Vertex",
		CurrentVersion: currentVersion,
		LatestVersion:  currentVersion,
		UpToDate:       true,
	}

	if currentVersion == "dev" {
		logger.Log("skipping update in 'dev' version").Print()
		return update, nil, nil
	}

	// remove previous old version if it exists.
	err := os.Remove("vertex-old")
	if err != nil && !os.IsNotExist(err) {
		return nil, nil, err
	}

	client := github.NewClient(nil)

	// get the latest release
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "vertex-center", "vertex")
	if err != nil {
		return nil, nil, err
	}

	// check if the version is different
	releaseVersion := *release.TagName
	releaseVersion = strings.TrimPrefix(releaseVersion, "v")

	update.LatestVersion = releaseVersion

	if currentVersion == releaseVersion {
		logger.Log("vertex is already up-to-date").Print()
		return update, nil, nil
	}

	logger.Log("a new release for Vertex is available").
		AddKeyValue("current", currentVersion).
		AddKeyValue("release", releaseVersion).
		Print()

	update.UpToDate = false
	return update, release, nil
}

func (s UpdateService) InstallVertexUpdate(currentVersion string) error {
	_, release, err := s.CheckForVertexUpdate(currentVersion)
	if err != nil {
		return err
	}

	err = storage.DownloadGithubRelease(release, storage.PathUpdates)
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

	logger.Warn("a new Vertex update has been installed. please restart Vertex.").Print()
	return nil
}
