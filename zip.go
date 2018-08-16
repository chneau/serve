package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func zipit(source string, target io.Writer) error {
	archive := zip.NewWriter(target)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == "upload" {
			return filepath.SkipDir
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			path := strings.TrimPrefix(path, source)
			if len(path) > 0 && (path[0] == '/' || path[0] == '\\') {
				path = path[1:]
			}
			if len(path) == 0 {
				return nil
			}
			header.Name = path
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
