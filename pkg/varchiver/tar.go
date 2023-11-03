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
	"strings"
)

var (
	ErrZipSlipAttack = errors.New("security: paths must be local")
)

func Unzip(src string, dest string) error {
	if zipSlipAttack(src) || zipSlipAttack(dest) {
		return ErrZipSlipAttack
	}

	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	for _, header := range reader.File {
		if zipSlipAttack(header.Name) {
			return ErrZipSlipAttack
		}

		p := path.Join(dest, header.Name)

		if zipSlipAttack(p) {
			return ErrZipSlipAttack
		}

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
	if zipSlipAttack(src) || zipSlipAttack(dest) {
		return ErrZipSlipAttack
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
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		if zipSlipAttack(header.Name) {
			return ErrZipSlipAttack
		}

		p := path.Join(dest, header.Name)

		if zipSlipAttack(p) {
			return ErrZipSlipAttack
		}

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

// CWE-22: Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')
func zipSlipAttack(path string) bool {
	return strings.Contains(path, "..")
}
