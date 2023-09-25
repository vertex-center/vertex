package varchiver

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

var (
	ErrMustBeLocal = errors.New("security: src must be a local file")
)

func Unzip(src string, dest string) error {
	if !filepath.IsLocal(src) || !filepath.IsLocal(dest) {
		return ErrMustBeLocal
	}

	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	for _, header := range reader.File {
		if !filepath.IsLocal(header.Name) {
			return ErrMustBeLocal
		}

		p := path.Join(dest, header.Name)

		if header.FileInfo().IsDir() {
			err = os.MkdirAll(p, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(path.Dir(p), os.ModePerm)
			if err != nil {
				return err
			}

			file, err := os.Create(p)
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

			err = os.Chmod(p, 0755)
			if err != nil {
				return err
			}

			file.Close()
		}
	}

	return nil
}

// Untar a tarball to a destination. src is the path to
// the tarball, and dest is the path to the destination directory.
func Untar(src string, dest string) error {
	if !filepath.IsLocal(src) || !filepath.IsLocal(dest) {
		return ErrMustBeLocal
	}

	archive, err := os.Open(src)
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

		if !filepath.IsLocal(header.Name) {
			return ErrMustBeLocal
		}

		p := path.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(p, os.ModePerm)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			err := os.MkdirAll(path.Dir(p), os.ModePerm)
			if err != nil {
				return err
			}

			file, err := os.Create(p)
			if err != nil {
				return err
			}

			_, err = io.Copy(file, reader)
			if err != nil {
				return err
			}

			err = os.Chmod(p, 0755)
			if err != nil {
				return err
			}

			file.Close()
		default:
			return fmt.Errorf("unknown flag type (%b) for file '%s'", header.Typeflag, header.Name)
		}
	}

	return nil
}
