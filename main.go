package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/howeyc/gopass"
	"github.com/klauspost/compress/zip"
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
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Value: "8888"},
				&cli.StringFlag{Name: "secret", Aliases: []string{"s"}},
			},
			Action: func(c *cli.Context) error {
				secret := c.String("secret")
				dir := c.Args().First()
				if dir == "" {
					dir = "."
				}
				askWhile("Secret: ", &secret)
				basePath := filepath.Dir(dir)
				total := uint64(0)
				filepath.Walk(dir, func(filePath string, fi os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					relativeFilePath, err := filepath.Rel(basePath, filePath)
					if err != nil {
						return err
					}
					_ = relativeFilePath
					if !fi.Mode().IsRegular() {
						return nil
					}
					header, err := zip.FileInfoHeader(fi)
					if err != nil {
						return err
					}
					total += header.UncompressedSize64
					return nil
				})
				log.Println("total", total)
				return nil
			},
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
			},
			Action: func(c *cli.Context) error {
				username := c.String("username")
				password := c.String("password")
				dir := c.Args().First()
				if dir == "" {
					dir = "."
				}
				dir, _ = filepath.Abs(dir)
				askWhile("Username: ", &username)
				askWhile("Password: ", &password)
				return web(dir, c.String("port"), password, username)
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
