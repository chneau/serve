package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mdp/qrterminal/v3"
	"github.com/samber/lo"
)

func GetOutboundIP() net.IP {
	conn := lo.Must(net.Dial("udp", "8.8.8.8:80"))
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func printIP(port string) {
	address := fmt.Sprintf("http://%s:%s/", GetOutboundIP(), port)
	log.Println("Listening on", address)
	qrterminal.GenerateHalfBlock(address, qrterminal.L, os.Stdout)
}
