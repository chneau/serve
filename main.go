package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	app := cli.NewApp()
	app.Name = "serve"
	app.Usage = "serve a file for direct download"
	app.Commands = []*cli.Command{
		{
			Name:      "send",
			Aliases:   []string{"s"},
			ArgsUsage: "[port]",
			Usage:     "send a folder using `serve receive`",
			Action:    sendAction,
		},
		{
			Name:      "receive",
			Aliases:   []string{"r"},
			ArgsUsage: "[ip:port]",
			Usage:     "send a folder",
			Action:    receiveAction,
		},
		{
			Name:      "web",
			Aliases:   []string{"w"},
			ArgsUsage: "[path]",
			Usage:     "web server to download or upload files",
			Flags:     []cli.Flag{&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Value: "8888"}},
			Action:    webAction,
		},
	}
	app.ArgsUsage = "[path]"
	app.Action = sendFileAction
	lo.Must0(app.Run(os.Args))
}
