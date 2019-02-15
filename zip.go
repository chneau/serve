package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"
)

func zipit(source string, target io.Writer) error {
	zw := zip.NewWriter(target)
	defer zw.Close()
	return filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
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
		header.Name = strings.TrimPrefix(strings.Replace(file, source, "", -1), string(filepath.Separator))
		header.Method = zip.Store
		wh, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(wh, f)
		return err
	})
}
