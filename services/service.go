package services

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/google/uuid"
)

type EnvVariable struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Secret      bool   `json:"secret,omitempty"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description"`
}

type Service struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Repository   string        `json:"repository"`
	Description  string        `json:"description"`
	EnvVariables []EnvVariable `json:"environment,omitempty"`
}

func (s Service) Install() (*Instance, error) {
	if strings.HasPrefix(s.Repository, "github") {
		client := github.NewClient(nil)

		split := strings.Split(s.Repository, "/")

		owner := split[1]
		repo := split[2]

		release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to retrieve the latest github release for %s", s.Repository))
		}

		platform := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)

		serviceUUID := uuid.New()

		for _, asset := range release.Assets {
			if strings.Contains(*asset.Name, platform) {
				basePath := path.Join("servers", serviceUUID.String())
				archivePath := path.Join(basePath, fmt.Sprintf("%s.tar.gz", s.ID))

				err := downloadFile(*asset.BrowserDownloadURL, basePath, archivePath)
				if err != nil {
					return nil, err
				}

				err = untarFile(basePath, archivePath)
				if err != nil {
					return nil, err
				}

				err = os.Remove(archivePath)
				if err != nil {
					return nil, err
				}

				break
			}
		}

		instance, err := Instantiate(serviceUUID)
		if err != nil {
			return nil, err
		}

		return instance, nil
	}

	return nil, errors.New("this repository is not supported")
}

func downloadFile(url string, basePath string, archivePath string) error {
	err := os.Mkdir(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func untarFile(basePath string, archivePath string) error {
	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	stream, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}
	defer stream.Close()

	reader := tar.NewReader(stream)

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		filepath := path.Join(basePath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filepath, os.ModePerm)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			err := os.MkdirAll(path.Dir(filepath), os.ModePerm)
			if err != nil {
				return err
			}

			file, err := os.Create(filepath)
			if err != nil {
				return err
			}

			_, err = io.Copy(file, reader)
			if err != nil {
				return err
			}

			err = os.Chmod(filepath, 0755)
			if err != nil {
				return err
			}

			file.Close()
		default:
			return errors.New(fmt.Sprintf("unknown flag type (%s) for file '%s'", header.Typeflag, header.Name))
		}
	}

	return nil
}
