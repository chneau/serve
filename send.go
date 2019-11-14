package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func getBytes(key interface{}) []byte {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(key)
	ce(err, "getBytes")
	return buf.Bytes()
}

func getFiles(dir string) map[string]uint64 {
	files := map[string]uint64{}
	err := filepath.Walk(dir, func(filePath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		files[filePath] = uint64(fi.Size())
		return nil
	})
	ce(err, "getFiles")
	return files
}

func sendAction(c *cli.Context) error {
	files := getFiles(".")
	log.Println(files)
	return nil
}
