package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/howeyc/gopass"
)

var path string
var port string
var password string
var username string
var noauth bool

func init() {
	gracefulExit()
	flag.StringVar(&path, "path", ".", "path to directory to serve")
	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.StringVar(&username, "usr", "", "username for auth")
	flag.StringVar(&password, "pwd", "", "password for auth")
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

func parse() {
	flag.Parse()
	if noauth == false {
		if username == "" {
			fmt.Print("Username: ")
			b, err := gopass.GetPasswdMasked()
			ce(err, "gopass.GetPasswd")
			username = string(b)
		}
		if password == "" {
			fmt.Print("Password: ")
			b, err := gopass.GetPasswdMasked()
			ce(err, "gopass.GetPasswd")
			password = string(b)
		}
	}
}

func main() {
	parse()
	r := gin.Default()
	if password != "" && username != "" {
		authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
			username: password,
		}))
		authorized.StaticFS("/", http.Dir(path))
	} else {
		r.StaticFS("/", http.Dir(path))
	}
	log.Println("Serving files from ", path)
	log.Printf("Listening on http://localhost:%s/\n", port)
	err := r.Run(":" + port)
	ce(err, "http.ListenAndServe")
}
