package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zip"
	"github.com/urfave/cli"
)

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

func sendAction(c *cli.Context) error {
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
}
