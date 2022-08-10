package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func getFiles(dir string) map[string]uint64 {
	files := map[string]uint64{}
	lo.Must0(filepath.Walk(dir, func(filePath string, fi os.FileInfo, _ error) error {
		if !fi.Mode().IsRegular() {
			return nil
		}
		files[filePath] = uint64(fi.Size())
		return nil
	}))
	return files
}

func sendAction(c *cli.Context) error {
	files := getFiles(".")
	port := c.String("port")
	r := gin.Default()
	r.GET("/files", func(c *gin.Context) {
		b := lo.Must(io.ReadAll(c.Request.Body))
		if len(b) == 0 {
			c.JSON(200, files)
			return
		}
		filename := string(b)
		if _, exist := files[filename]; !exist {
			_ = c.Error(errors.New("the file requested does not exist"))
			return
		}
		f := lo.Must(os.Open(filename))
		defer f.Close()
		lo.Must(io.Copy(c.Writer, f))
	})
	r.GET("/end", func(c *gin.Context) {
		c.Writer.Flush()
		os.Exit(0)
	})
	printIP(port)
	return r.Run(":" + port)
}
