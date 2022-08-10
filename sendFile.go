package main

import (
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func sendFileAction(c *cli.Context) error {
	port := "8888"
	fileName := c.Args().First()
	fileInfo := lo.Must(os.Stat(fileName))
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		header := c.Writer.Header()
		header["Content-Type"] = []string{"application/octet-stream"}
		header["Content-Length"] = []string{strconv.FormatInt(fileInfo.Size(), 10)}
		header["Content-Disposition"] = []string{"attachment; filename=" + fileInfo.Name()}
		f := lo.Must(os.Open(fileName))
		defer f.Close()
		lo.Must(io.Copy(c.Writer, f))
	})
	printIP(port)
	return r.Run(":" + port)
}
