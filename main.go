package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/howeyc/gopass"
)

const (
	empty = ""
	route = "/"
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
	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.StringVar(&username, "usr", empty, "username for auth")
	flag.StringVar(&password, "pwd", empty, "password for auth")
	flag.BoolVar(&noauth, "noauth", false, "do not ask for auth")
	gin.SetMode(gin.ReleaseMode)
}

// checkError
func ce(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

// Handle the ctrl+c with some grace
func gracefulExit() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		os.Exit(0)
	}()
}

// Ask something to hide secretly to the user
func askWhile(prompt string) string {
	res := empty
	for res == empty {
		b, err := gopass.GetPasswdPrompt(prompt, true, os.Stdin, os.Stdout)
		ce(err, "gopass.GetPasswdPrompt")
		res = string(b)
	}
	return res
}

// Parse flags
func parse() {
	flag.Parse()
	if noauth == false {
		username = askWhile("Username: ")
		password = askWhile("Password: ")
	}
}

func main() {
	parse()
	r := gin.Default()
	if password != empty && username != empty {
		r.Group(route, gin.BasicAuth(gin.Accounts{username: password})).
			StaticFS(route, http.Dir(path))
	} else {
		r.StaticFS(route, http.Dir(path))
	}
	log.Println("Serving files from ", path)
	log.Printf("Listening on http://localhost:%s/\n", port)
	err := r.Run(":" + port)
	ce(err, "http.ListenAndServe")
}
