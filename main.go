package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/chneau/serve/pkg/statik"
	"github.com/gin-gonic/gin"
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
			Name:    "web",
			Aliases: []string{"w"},
			Usage:   "web page serving",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dir", Aliases: []string{"d"}, Value: "."},
				&cli.StringFlag{Name: "port", Aliases: []string{"n"}, Value: "8888"},
				&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
				&cli.StringFlag{Name: "username", Aliases: []string{"u"}},
				&cli.BoolFlag{Name: "auth", Aliases: []string{"a"}, Value: true},
			},
			Action: func(c *cli.Context) error {
				auth := c.Bool("auth")
				username := c.String("username")
				password := c.String("password")
				dir := c.String("dir")
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

func web(dir, port, password, username string, auth bool) error {
	gin.SetMode(gin.ReleaseMode)
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	}
	html, err := statik.Asset("public/index.html")
	ce(err, `statik.Asset("public/index.html")`)
	dcss, err := statik.Asset("public/dropzone.css")
	ce(err, `statik.Asset("public/dropzone.css")`)
	djs, err := statik.Asset("public/dropzone.js")
	ce(err, `statik.Asset("public/dropzone.js")`)
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
		os.MkdirAll(dir+"/uploaded_files/"+fullPath[:len(fullPath)-len(file.Filename)], 0777)
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
		zipit(cleanedPath, c.Writer)
	})
	printIP(port)
	return r.Run(":" + port)
}
