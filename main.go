package main

import (
	"encoding/gob"
	"log"
	"net"
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
func askWhile(prompt string) string {
	res := ""
	for res == "" {
		b, err := gopass.GetPasswdPrompt(prompt, true, os.Stdin, os.Stdout)
		if err != nil {
			os.Exit(0)
		}
		res = string(b)
	}
	return res
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
			},
			Action: func(c *cli.Context) error {
				dir := c.Args().First()
				if dir == "" {
					dir = "."
				}
				dir, _ = filepath.Rel(".", dir)
				secret := askWhile("Secret: ")
				_ = secret
				basePath := filepath.Dir(dir)
				_ = basePath
				files := map[string]int{}
				err := filepath.Walk(dir, func(filePath string, fi os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !fi.Mode().IsRegular() {
						return nil
					}
					header, err := zip.FileInfoHeader(fi)
					if err != nil {
						return err
					}
					files[filePath] = int(header.UncompressedSize)
					return nil
				})
				if err != nil {
					return err
				}
				log.Println(files)
				return nil
			},
		},
		{
			Name:      "receive",
			Aliases:   []string{"r"},
			ArgsUsage: "[ip]",
			Usage:     "send a folder",
			Action: func(c *cli.Context) error {
				secret := askWhile("Secret: ")
				conn, err := net.Dial("tcp4", ":8888")
				if err != nil {
					return err
				}
				_, err = conn.Write([]byte(secret))
				if err != nil {
					return err
				}
				files := map[string]int{}
				err = gob.NewDecoder(conn).Decode(&files)
				if err != nil {
					return err
				}
				log.Println(files)
				return nil
			},
		},
		{
			Name:      "web",
			Aliases:   []string{"w"},
			ArgsUsage: "[path]",
			Usage:     "web page serving",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Value: "8888"},
			},
			Action: func(c *cli.Context) error {
				dir := c.Args().First()
				if dir == "" {
					dir = "."
				}
				dir, _ = filepath.Abs(dir)
				username := askWhile("Username: ")
				password := askWhile("Password: ")
				return web(dir, c.String("port"), password, username)
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
