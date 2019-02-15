package main

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/pgzip"
)

func tgzit(source string, target io.Writer) error {
	gzw, _ := pgzip.NewWriterLevel(target, pgzip.HuffmanOnly)
	defer gzw.Close()
	tw := tar.NewWriter(gzw)
	defer tw.Close()
	return filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, source, "", -1), string(filepath.Separator))
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(tw, f)
		return nil
	})
}
