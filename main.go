package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/howeyc/gopass"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetPrefix("[SRV] ")
	log.SetFlags(log.LstdFlags)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		println()
		os.Exit(0)
	}()
}

func main() {
	app := cli.NewApp()
	app.Name = "serve"
	app.Usage = "serve files from or to another computer"
	app.Version = "0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:    "send",
			Aliases: []string{"s"},
			Usage:   "send a folder",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Value: "8888"},
			},
			Action: sendAction,
		},
		{
			Name:      "receive",
			Aliases:   []string{"r"},
			ArgsUsage: "[ip:port]",
			Usage:     "send a folder",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "concurrence", Aliases: []string{"c"}, Value: 100},
			},
			Action: receiveAction,
		},
		{
			Name:      "web",
			Aliases:   []string{"w"},
			ArgsUsage: "[path]",
			Usage:     "web page serving",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Value: "8888"},
			},
			Action: webAction,
		},
	}
	app.Action = func(c *cli.Context) error {
		log.Println("Welcome to serve!")
		return nil
	}
	lo.Must0(app.Run(os.Args))
}

// Ask something to hide secretly to the user
func askWhile(prompt string) string {
	res := ""
	for res == "" {
		b := lo.Must(gopass.GetPasswdPrompt(prompt, true, os.Stdin, os.Stdout))
		res = string(b)
	}
	return res
}
