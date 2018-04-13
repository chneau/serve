package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/howeyc/gopass"
)

var (
	path     string
	port     string
	password string
	username string
	noauth   bool
)

func init() {
	gracefulExit()
	flag.StringVar(&path, "path", ".", "path to directory to serve")
	flag.StringVar(&port, "port", "8888", "port to listen on")
	flag.StringVar(&username, "usr", "", "username for auth")
	flag.StringVar(&password, "pwd", "", "password for auth")
	flag.BoolVar(&noauth, "noauth", false, "do not ask for auth")
	gin.SetMode(gin.ReleaseMode)
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

// Parse flags and ask username / password
func parse() {
	flag.Parse()
	if noauth == false {
		askWhile("Username: ", &username)
		askWhile("Password: ", &password)
	}
}

func serveGroup(r *gin.Engine) *gin.RouterGroup {
	if password != "" && username != "" {
		return r.Group("/", gin.BasicAuth(gin.Accounts{username: password}))
	}
	return r.Group("/")
}

func serve() {
	r := gin.Default()
	r.Use(gin.Recovery())
	grp := serveGroup(r)
	grp.GET("/", func(c *gin.Context) { c.Redirect(307, "/serve") })
	grp.StaticFS("/serve", http.Dir(path))
	grp.GET("/upload", func(c *gin.Context) {
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
		ff.Close()
		f.Close()
		if written != file.Size {
			c.Status(406)
		}
		c.Status(201)
	})
	log.Println("Serving files from", path)
	hostname, err := os.Hostname()
	ce(err, "os.Hostname")
	log.Printf("Listening on http://%[1]s:%[2]s/serve\n\thttp://localhost:%[2]s/serve\n", hostname, port)
	log.Printf("Listening on http://%[1]s:%[2]s/upload\n\thttp://localhost:%[2]s/upload\n", hostname, port)
	err = r.Run(":" + port)
	ce(err, "http.ListenAndServe")
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
	parse()
	serve()
}
