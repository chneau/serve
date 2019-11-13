package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
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
	log.SetFlags(log.LstdFlags | log.Llongfile)
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

func getBytes(key interface{}) []byte {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(key)
	ce(err, "getBytes")
	return buf.Bytes()
}

func getFiles(dir string) map[string]uint64 {
	files := map[string]uint64{}
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
		files[filePath] = header.UncompressedSize64
		return nil
	})
	ce(err, "getFiles")
	return files
}

func uint64ToBytes(x uint64) []byte {
	bb := make([]byte, 8)
	binary.LittleEndian.PutUint64(bb, x)
	return bb
}

func uint64FromConn(conn net.Conn) (uint64, error) {
	b := make([]byte, 8)
	_, err := conn.Read(b)
	if err != nil {
		return 0, err
	}
	ce(err, "uint64FromConn")
	res := binary.LittleEndian.Uint64(b)
	return res, nil
}

func filesFromConn(conn net.Conn, size uint64) map[string]uint64 {
	log.Println(size)
	b := make([]byte, size)
	_, err := conn.Read(b)
	ce(err, "filesFromConn")
	res := map[string]uint64{}
	gob.NewDecoder(bytes.NewBuffer(b)).Decode(&res)
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
				dir, _ := filepath.Rel(".", c.Args().First())
				// secret := askWhile("Secret: ")
				// _ = secret
				files := getFiles(dir)
				log.Println("len(files)", len(files))
				b := getBytes(files)
				log.Println(len(b))
				listener, err := net.Listen("tcp4", ":"+c.String("port"))
				ce(err, "net.Listen")
				conn, err := listener.Accept()
				ce(err, "listener.Accept")
				_, err = conn.Write(uint64ToBytes(uint64(len(b))))
				ce(err, "conn.Write")
				_, err = conn.Write(b)
				ce(err, "conn.Write")
				for file, fsize := range files {
					ssize := uint64ToBytes(uint64(len(file)))
					// log.Println("ssize", len(file))
					_, err = conn.Write(ssize)
					ce(err, "conn.Write")
					_, err = conn.Write([]byte(file))
					ce(err, "conn.Write")
					_, err = conn.Write(uint64ToBytes(fsize))
					// log.Println("fsize", fsize)
					ce(err, "conn.Write")
					f, err := os.Open(file)
					ce(err, "os.Open")
					_, err = io.Copy(conn, f)
					ce(err, "io.Copy")
				}
				return nil
			},
		},
		{
			Name:      "receive",
			Aliases:   []string{"r"},
			ArgsUsage: "[ip]",
			Usage:     "send a folder",
			Action: func(c *cli.Context) error {
				// secret := askWhile("Secret: ")
				conn, err := net.Dial("tcp4", ":8888")
				if err != nil {
					return err
				}
				size, err := uint64FromConn(conn)
				if err != nil {
					return err
				}
				files := filesFromConn(conn, size)
				log.Println("len(files)", len(files))
				for len(files) > 0 {
					size, err := uint64FromConn(conn)
					if err != nil {
						return err
					}
					log.Println("ssize", size)
					b := make([]byte, size)
					_, err := conn.Read(b)
					ce(err, "conn.Read")
					log.Println(string(b))
					size, err = uint64FromConn(conn)
					if err != nil {
						return err
					}
					log.Println(size)
					b = make([]byte, size)
					_, err = conn.Read(b)
					ce(err, "conn.Read")
					delete(files, string(b))
				}
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
