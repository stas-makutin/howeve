package utils

import (
	"archive/zip"
	"io"
	"os"
	"strings"
)

// FileToArchive struct
type FileToArchive struct {
	Name, Path string
}

// ZipFilesToWriter func
func ZipFilesToWriter(w *zip.Writer, files []FileToArchive) error {
	for _, file := range files {
		err := func() error {
			src, err := os.Open(file.Path)
			if err != nil {
				return err
			}
			defer src.Close()

			dest, err := w.Create(file.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(dest, src); err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

// ZipFilesToFile func
func ZipFilesToFile(zipFile string, perm os.FileMode, files []FileToArchive) error {
	f, err := os.OpenFile(zipFile, os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	err = func() error {
		defer f.Close()
		zw := zip.NewWriter(f)
		err := ZipFilesToWriter(zw, files)
		errClose := zw.Close()
		if err != nil {
			return err
		}
		return errClose
	}()
	if err != nil {
		os.Remove(zipFile)
	}
	return err
}

// WriteStringln func
func WriteStringln(sb *strings.Builder, s string) {
	if sb.Len() > 0 {
		sb.WriteString(NewLine)
	}
	sb.WriteString(s)
}
