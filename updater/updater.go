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

func CheckForUpdates(currentVersion string) error {
	if currentVersion == "dev" {
		logger.Log("skipping update in 'dev' version")
		return nil
	}

	// remove previous old version if it exists.
	err := os.Remove("vertex-old")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	client := github.NewClient(nil)

	// get the latest release
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "vertex-center", "vertex")
	if err != nil {
		return err
	}

	// check if the version is different
	releaseVersion := *release.TagName
	if strings.HasPrefix(releaseVersion, "v") {
		releaseVersion = releaseVersion[1:]
	}

	if currentVersion == releaseVersion {
		logger.Log("vertex is already up-to-date")
		return nil
	}

	logger.Log(fmt.Sprintf("a new release is available ('%s'), currently using '%s'", releaseVersion, currentVersion))

	// download if it is newer
	err = storage.DownloadGithubRelease(release, storage.PathUpdates)
	if err != nil {
		return err
	}

	// replace the old executable with the new one
	err = os.Rename("vertex", "vertex-old")
	if err != nil {
		return err
	}

	err = os.Rename(path.Join(storage.PathUpdates, "vertex"), "vertex")
	if err != nil {
		return err
	}

	logger.Warn("A new update has been installed. Please restart Vertex.")
	os.Exit(0)
	return nil
}
