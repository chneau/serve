package main

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/klauspost/compress/zip"
	"github.com/samber/lo"
)

func zipit(inFilePath string, target io.Writer) {
	basePath := filepath.Dir(inFilePath)
	zw := zip.NewWriter(target)
	defer zw.Close()
	filepath.Walk(inFilePath, func(filePath string, fi os.FileInfo, _ error) error {
		relativeFilePath := lo.Must(filepath.Rel(basePath, filePath))
		if !fi.Mode().IsRegular() {
			return nil
		}
		header := lo.Must(zip.FileInfoHeader(fi))
		header.Name = path.Join(filepath.SplitList(relativeFilePath)...)
		header.Method = zip.Store
		wh := lo.Must(zw.CreateHeader(header))
		f := lo.Must(os.Open(filePath))
		defer f.Close()
		lo.Must(io.Copy(wh, f))
		return nil
	})
}
