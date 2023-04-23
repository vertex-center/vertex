package client

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/vertex-center/vertex/logger"
	"github.com/vertex-center/vertex/storage"
)

func Setup() error {
	err := os.Mkdir(storage.PathClient, os.ModePerm)
	if os.IsExist(err) {
		// The client is already setup.
		return nil
	}
	if err != nil {
		return err
	}

	client := github.NewClient(nil)

	owner := "vertex-center"
	repo := "vertex-webui"

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		return fmt.Errorf("failed to retrieve the latest github release for %s: %v", repo, err)
	}

	logger.Log("downloading vertex-webui client...").Print()

	for _, asset := range release.Assets {
		if strings.Contains(*asset.Name, "vertex-webui") {
			err := download(*asset.BrowserDownloadURL)
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

	return nil
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
