package main

import (
	"os"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "serve"
	app.Usage = "serve files from or to another computer"
	app.Version = "0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:      "send",
			Aliases:   []string{"s"},
			ArgsUsage: "[port]",
			Usage:     "send a folder",
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
			Usage:     "web to download or upload files",
			Flags:     []cli.Flag{&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Value: "8888"}},
			Action:    webAction,
		},
	}
	lo.Must0(app.Run(os.Args))
}
