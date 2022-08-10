package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/ioprogress"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func sendFileAction(c *cli.Context) error {
	port := "8888"
	fileName := c.Args().First()
	fileInfo := lo.Must(os.Stat(fileName))
	r := gin.Default()
	bar := ioprogress.DrawTextFormatBar(40)
	drawFunc := func(progress, total int64) string {
		return fmt.Sprintf("%s %s\r", bar(progress, total), ioprogress.DrawTextFormatBytes(progress, total))
	}
	r.GET("/", func(c *gin.Context) {
		header := c.Writer.Header()
		header["Content-Type"] = []string{"application/octet-stream"}
		header["Content-Length"] = []string{strconv.FormatInt(fileInfo.Size(), 10)}
		header["Content-Disposition"] = []string{"attachment; filename=" + fileInfo.Name()}
		f := lo.Must(os.Open(fileName))
		defer f.Close()
		lo.Must(io.Copy(c.Writer, &ioprogress.Reader{
			Reader:       f,
			Size:         fileInfo.Size(),
			DrawFunc:     ioprogress.DrawTerminalf(os.Stdout, drawFunc),
			DrawInterval: time.Millisecond * 100,
		}))
		defer os.Exit(0)
	})
	printIP(port)
	return r.Run(":" + port)
}
