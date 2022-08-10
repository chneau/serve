package main

import (
	_ "embed"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

//go:embed public/index.html
var html []byte

//go:embed public/dropzone.css
var dcss []byte

//go:embed public/dropzone.js
var djs []byte

func web(dir, port string) {
	r := gin.Default()
	r.StaticFS("/serve", http.Dir(dir))
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charsed=ute-8", html)
	})
	r.GET("/dropzone.js", func(c *gin.Context) {
		c.Data(200, "text/javascript; charsed=ute-8", djs)
	})
	r.GET("/dropzone.css", func(c *gin.Context) {
		c.Data(200, "text/css; charsed=ute-8", dcss)
	})
	r.POST("/upload", func(c *gin.Context) {
		file := lo.Must(c.FormFile("file"))
		fullPath := c.PostForm("fullPath")
		lo.Must0(os.MkdirAll(dir+"/uploaded_files/"+fullPath[:len(fullPath)-len(file.Filename)], 0777))
		f := lo.Must(os.OpenFile(dir+"/uploaded_files/"+fullPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666))
		defer f.Close()
		ff := lo.Must(file.Open())
		defer ff.Close()
		written := lo.Must(io.Copy(f, ff))
		if written != file.Size {
			c.Status(406)
		}
		c.Status(201)
	})
	r.GET("/zip/*path", func(c *gin.Context) {
		p := c.Param("path")
		cleanedPath := filepath.Clean(dir + p)
		header := c.Writer.Header()
		header["Content-Disposition"] = []string{"attachment; filename= " + filepath.Base(cleanedPath) + ".zip"}
		zipit(cleanedPath, c.Writer)
	})
	printIP(port)
	lo.Must0(r.Run(":" + port))
}

func webAction(c *cli.Context) error {
	dir := c.Args().First()
	if dir == "" {
		dir = "."
	}
	dir = lo.Must(filepath.Abs(dir))
	web(dir, c.String("port"))
	return nil
}
