package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	_ "github.com/chneau/serve/pkg/statik"
	"github.com/gin-gonic/gin"
	"github.com/howeyc/gopass"
	"github.com/rakyll/statik/fs"
)

var (
	path     string
	port     string
	password string
	username string
	noauth   bool
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	}
	gracefulExit()
	flag.StringVar(&path, "path", ".", "path to directory to serve")
	flag.StringVar(&port, "port", "8888", "port to listen on")
	flag.StringVar(&username, "usr", "", "username for auth")
	flag.StringVar(&password, "pwd", "", "password for auth")
	flag.BoolVar(&noauth, "noauth", false, "do not ask for auth")
	log.SetPrefix("[SRV] ")
	log.SetFlags(log.LstdFlags)
}

// checkError
func ce(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}

// Ask something to hide secretly to the user
func askWhile(prompt string, res *string) {
	for *res == "" {
		b, err := gopass.GetPasswdPrompt(prompt, true, os.Stdin, os.Stdout)
		ce(err, "gopass.GetPasswdPrompt")
		*res = string(b)
	}
}

func serveGroup(r *gin.Engine) *gin.RouterGroup {
	if password != "" && username != "" {
		return r.Group("/", gin.BasicAuth(gin.Accounts{username: password}))
	}
	return r.Group("/")
}

func gracefulExit() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		os.Exit(0)
	}()
}

func main() {

	fs, err := fs.New()
	ce(err, "fs.New()")

	f, err := fs.Open("/")
	ce(err, `fs.Open("/index.html")`)

	html, err := ioutil.ReadAll(f)
	ce(err, "ioutil.ReadAll(f)")

	f.Close()

	flag.Parse()
	if noauth == false {
		askWhile("Username: ", &username)
		askWhile("Password: ", &password)
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	grp := serveGroup(r)
	grp.StaticFS("/serve", http.Dir(path))
	grp.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charsed=ute-8", html)
	})
	grp.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		ce(err, "c.FormFile")
		fullPath := c.PostForm("fullPath")
		os.MkdirAll(path+"/upload/"+fullPath[:len(fullPath)-len(file.Filename)], 0777)
		f, err := os.OpenFile(path+"/upload/"+fullPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
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
	grp.GET("/zip", func(c *gin.Context) {
		zipit(path,c.Writer)
	})
	hostname, err := os.Hostname()
	ce(err, "os.Hostname")
	log.Printf("Listening on http://%[1]s:%[2]s/ , http://localhost:%[2]s/\n", hostname, port)
	err = r.Run(":" + port)
	ce(err, "http.ListenAndServe")
}
