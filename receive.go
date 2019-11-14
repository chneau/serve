package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"net"

	"github.com/urfave/cli"
)

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

func receiveAction(c *cli.Context) error {
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
		_, err = conn.Read(b)
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
}
