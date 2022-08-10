package main

import (
	"log"
	"net"

	"github.com/samber/lo"
)

func GetOutboundIP() net.IP {
	conn := lo.Must(net.Dial("udp", "8.8.8.8:80"))
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func printIP(port string) {
	log.Printf("Listening on http://%s:%s/", GetOutboundIP(), port)
}
