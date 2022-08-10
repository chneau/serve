package main

import (
	"log"
	"net"

	"github.com/samber/lo"
)

func printIP(port string) {
	ifaces := lo.Must(net.Interfaces())
	for _, i := range ifaces {
		addrs := lo.Must(i.Addrs())
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
