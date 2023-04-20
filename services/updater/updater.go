package updater

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/storage"
)

var logger = console.New("updater")

type Update struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	UpToDate       bool   `json:"up_to_date"`
}

func CheckForUpdates(currentVertexVersion string) ([]Update, error) {
	var updates []Update

	vertexUpdate, _, err := CheckForVertexUpdate(currentVertexVersion)
	if err != nil {
		return nil, err
	} else {
		updates = append(updates, *vertexUpdate)
	}

	return updates, nil
}

func CheckForVertexUpdate(currentVersion string) (*Update, *github.RepositoryRelease, error) {
	update := &Update{
		Id:             "vertex",
		Name:           "Vertex",
		CurrentVersion: currentVersion,
		LatestVersion:  currentVersion,
		UpToDate:       true,
	}

	if currentVersion == "dev" {
		logger.Log("skipping update in 'dev' version")
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
		logger.Log("vertex is already up-to-date")
		return update, nil, nil
	}

	logger.Log(fmt.Sprintf("a new release for Vertex is available ('%s'), currently using '%s'", releaseVersion, currentVersion))

	update.UpToDate = false
	return update, release, nil
}

func InstallVertexUpdate(currentVersion string) error {
	_, release, err := CheckForVertexUpdate(currentVersion)
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

	logger.Warn("A new Vertex update has been installed. Please restart Vertex.")
	return nil
}
