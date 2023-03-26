package client

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
	"github.com/vertex-center/vertex-core-golang/console"
)

var logger = console.New("client")

func Setup() error {
	err := os.Mkdir("clients", os.ModePerm)
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
		return errors.New(fmt.Sprintf("failed to retrieve the latest github release for %s: %v", repo, err))
	}

	logger.Log("Downloading vertex-webui client...")

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

			err = os.Remove(path.Join("clients", "temp.zip"))
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

	file, err := os.Create(path.Join("clients", "temp.zip"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func unarchive() error {
	reader, err := zip.OpenReader(path.Join("clients", "temp.zip"))
	if err != nil {
		return err
	}

	for _, header := range reader.File {
		filepath := path.Join("clients", header.Name)

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
