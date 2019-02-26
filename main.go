package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	_ "github.com/chneau/serve/pkg/statik"
	"github.com/gin-gonic/gin"
	"github.com/howeyc/gopass"
	"github.com/rakyll/statik/fs"
)

var (
	pathDir  string
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
	flag.StringVar(&pathDir, "path", ".", "path to directory to serve")
	flag.StringVar(&port, "port", "8888", "port to listen on")
	flag.StringVar(&username, "usr", "", "username for auth")
	flag.StringVar(&password, "pwd", "", "password for auth")
	flag.BoolVar(&noauth, "noauth", false, "do not ask for auth")
	log.SetPrefix("[SRV] ")
	log.SetFlags(log.LstdFlags)
	pathDir, _ = filepath.Abs(pathDir)
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
		println()
		os.Exit(0)
	}()
}

func printIP(port string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() == nil {
				continue
			}
			log.Printf("Listening on (%s) http://%s:%s/", i.Name, ip, port)
		}
	}
}

func main() {

	fs, err := fs.New()
	ce(err, "fs.New()")

	f, err := fs.Open("/index.html")
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
	grp.StaticFS("/serve", http.Dir(pathDir))
	grp.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charsed=ute-8", html)
	})
	grp.GET("/dropzone.js", func(c *gin.Context) {
		res, err := http.DefaultClient.Get("https://cdn.jsdelivr.net/npm/dropzone@5.5.1/dist/min/dropzone.min.js")
		if err != nil {
			c.Error(err)
		}
		c.Header("Content-Type", "text/javascript")
		io.Copy(c.Writer, res.Body)
	})
	grp.GET("/dropzone.css", func(c *gin.Context) {
		res, err := http.DefaultClient.Get("https://cdn.jsdelivr.net/npm/dropzone@5.5.1/dist/min/dropzone.min.css")
		if err != nil {
			c.Error(err)
		}
		c.Header("Content-Type", "text/css")
		io.Copy(c.Writer, res.Body)
	})
	grp.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		ce(err, "c.FormFile")
		fullPath := c.PostForm("fullPath")
		os.MkdirAll(pathDir+"/uploaded_files/"+fullPath[:len(fullPath)-len(file.Filename)], 0777)
		f, err := os.OpenFile(pathDir+"/uploaded_files/"+fullPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
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
		cleanedPath := filepath.Clean(pathDir + p)
		header := c.Writer.Header()
		header["Content-Disposition"] = []string{"attachment; filename= " + filepath.Base(cleanedPath) + ".zip"}
		zipit(cleanedPath, c.Writer)
	})
	printIP(port)
	err = r.Run(":" + port)
	ce(err, "http.ListenAndServe")
}
