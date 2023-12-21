package updates

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/varchiver"
	"github.com/vertex-center/vlog"
)

type VertexClientUpdater struct {
	dir string
}

func NewVertexClientUpdater(dir string) VertexClientUpdater {
	return VertexClientUpdater{
		dir: dir,
	}
}

func (u VertexClientUpdater) CurrentVersion() (string, error) {
	version, err := os.ReadFile(path.Join(u.dir, "dist", "version.txt"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(version)), nil
}

func (u VertexClientUpdater) Install(tag string) error {
	client := github.NewClient(nil)

	log.Info("installing vertex client", vlog.String("tag", tag))

	release, res, err := client.Repositories.GetReleaseByTag(context.Background(), "vertex-center", "client", tag)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	for _, asset := range release.Assets {
		if strings.Contains(*asset.Name, "client") {
			return install(u.dir, *asset.BrowserDownloadURL)
		}
	}

	return nil
}

func (u VertexClientUpdater) IsInstalled() bool {
	_, err := os.Stat(path.Join(u.dir, "dist"))
	return err == nil
}

func (u VertexClientUpdater) ID() string {
	return "vertex_client"
}

func install(dir string, releaseUrl string) error {
	tempPath := path.Join(dir, "temp.zip")

	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	err = download(tempPath, releaseUrl)
	if err != nil {
		return err
	}

	err = varchiver.Unzip(tempPath, dir)
	if err != nil {
		return err
	}

	err = os.Remove(tempPath)
	if err != nil {
		return err
	}

	return applyConfig()
}

func download(dir string, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(dir)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func applyConfig() error {
	cfg := "window.api_urls = {\n"
	// Only for the non-kernel apps
	for name, u := range config.Current.Urls {
		name = strings.ReplaceAll(name, "-", "_")
		cfg += fmt.Sprintf("\t%s: '%s',\n", name, u)
	}
	cfg += "};\n"
	return os.WriteFile(path.Join(storage.FSPath, "client", "dist", "config.js"), []byte(cfg), os.ModePerm)
}
