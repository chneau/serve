package main

import (
	_ "embed"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

//go:embed public/index.html
var html []byte

//go:embed public/dropzone.css
var dcss []byte

//go:embed public/dropzone.js
var djs []byte

func web(dir, port, password, username string) error {
	gin.SetMode(gin.ReleaseMode)
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	}
	r := gin.Default()
	r.Use(gin.Recovery())
	opts := []gin.HandlerFunc{}
	if password != "" && username != "" {
		opts = append(opts, gin.BasicAuth(gin.Accounts{username: password}))
	}
	grp := r.Group("/", opts...)
	grp.StaticFS("/serve", http.Dir(dir))
	grp.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charsed=ute-8", html)
	})
	grp.GET("/dropzone.js", func(c *gin.Context) {
		c.Data(200, "text/javascript; charsed=ute-8", djs)
	})
	grp.GET("/dropzone.css", func(c *gin.Context) {
		c.Data(200, "text/css; charsed=ute-8", dcss)
	})
	grp.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		ce(err, "c.FormFile")
		fullPath := c.PostForm("fullPath")
		err = os.MkdirAll(dir+"/uploaded_files/"+fullPath[:len(fullPath)-len(file.Filename)], 0777)
		ce(err, "os.MkdirAll")
		f, err := os.OpenFile(dir+"/uploaded_files/"+fullPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
		ce(err, "os.OpenFile")
		ff, err := file.Open()
		ce(err, "file.Open")
		written, err := io.Copy(f, ff)
		ce(err, "io.Copy")
		ce(ff.Close(), "ff.Close()")
		ce(f.Close(), "f.Close()")
		if written != file.Size {
			c.Status(406)
		}
		c.Status(201)
	})
	grp.GET("/zip/*path", func(c *gin.Context) {
		p := c.Param("path")
		cleanedPath := filepath.Clean(dir + p)
		header := c.Writer.Header()
		header["Content-Disposition"] = []string{"attachment; filename= " + filepath.Base(cleanedPath) + ".zip"}
		err := zipit(cleanedPath, c.Writer)
		ce(err, "zipit")
	})
	printIP(port)
	return r.Run(":" + port)
}

func webAction(c *cli.Context) error {
	dir := c.Args().First()
	if dir == "" {
		dir = "."
	}
	dir, _ = filepath.Abs(dir)
	username := askWhile("Username: ")
	password := askWhile("Password: ")
	return web(dir, c.String("port"), password, username)
}
