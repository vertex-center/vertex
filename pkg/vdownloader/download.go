package vdownloader

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vlog"
)

// Download downloads a file from a Addr. It creates the
// directory if it doesn't exist.
func Download(url string, dir string, filename string) error {
	log.Info("downloading",
		vlog.String("url", url),
	)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(path.Join(dir, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}
