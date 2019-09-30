package main

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/klauspost/compress/zip"
)

func zipit(inFilePath string, target io.Writer) error {
	basePath := filepath.Dir(inFilePath)
	zw := zip.NewWriter(target)
	defer zw.Close()
	return filepath.Walk(inFilePath, func(filePath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relativeFilePath, err := filepath.Rel(basePath, filePath)
		if err != nil {
			return err
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		header.Name = path.Join(filepath.SplitList(relativeFilePath)...)
		header.Method = zip.Store
		wh, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(wh, f)
		return err
	})
}
