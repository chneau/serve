package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

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
	port := c.String("port")
	gin.SetMode(gin.ReleaseMode)
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	}
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/files", func(c *gin.Context) {
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			_ = c.Error(err)
			return
		}
		if len(b) == 0 {
			c.JSON(200, files)
			return
		}
		filename := string(b)
		if _, exist := files[filename]; !exist {
			_ = c.Error(errors.New("the file requested does not exist"))
			return
		}
		f, err := os.Open(filename)
		if err != nil {
			_ = c.Error(err)
			return
		}
		defer f.Close()
		_, err = io.Copy(c.Writer, f)
		if err != nil {
			_ = c.Error(err)
			return
		}
	})
	r.GET("/end", func(_ *gin.Context) { os.Exit(0) })
	printIP(port)
	return r.Run(":" + port)
}
