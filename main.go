package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

func init() {
	log.SetPrefix("[SRV] ")
	log.SetFlags(log.LstdFlags)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		println()
		os.Exit(0)
	}()
}

// checkError
func ce(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

// Ask something to hide secretly to the user
func askWhile(prompt string, res *string) {
	for *res == "" {
		b, err := gopass.GetPasswdPrompt(prompt, true, os.Stdin, os.Stdout)
		if err != nil {
			os.Exit(0)
		}
		*res = string(b)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "serve"
	app.Usage = "serve files from or to another computer"
	app.Version = "0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:      "send",
			Aliases:   []string{"s"},
			ArgsUsage: "[path]",
			Usage:     "send a folder",
			Flags:     []cli.Flag{},
		},
		{
			Name:      "web",
			Aliases:   []string{"w"},
			ArgsUsage: "[path]",
			Usage:     "web page serving",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "port", Aliases: []string{"n"}, Value: "8888"},
				&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
				&cli.StringFlag{Name: "username", Aliases: []string{"u"}},
				&cli.BoolFlag{Name: "auth", Aliases: []string{"a"}, Value: true},
			},
			Action: func(c *cli.Context) error {
				auth := c.Bool("auth")
				username := c.String("username")
				password := c.String("password")
				dir := c.Args().First()
				if dir == "" {
					dir = "."
				}
				dir, _ = filepath.Abs(dir)
				if auth {
					askWhile("Username: ", &username)
					askWhile("Password: ", &password)
				}
				web(dir, c.String("port"), password, username, auth)
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
