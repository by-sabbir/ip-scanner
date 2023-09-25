package main

import (
	"fmt"
	"net"

	"github.com/by-sabbir/ip-scanner/ping"
	"golang.org/x/exp/slog"
)

func main() {
	cidr := "192.168.0.0/30"

	ip, netmask, err := net.ParseCIDR(cidr)

	if err != nil {
		slog.Error("error parsing cidr: ", err)
	}

	fmt.Println("net: ", netmask)
	fmt.Println("ip: ", ip)

	nextIp := ip
	aliveIps := []net.IP{}
	for {
		nextIp = getNextIP(nextIp, 1)
		fmt.Println("next ip: ", nextIp)
		alive, err := ping.Ping(nextIp)

		if err != nil {
			slog.Error(err.Error())
		}

		if alive {
			aliveIps = append(aliveIps, nextIp)
		}

		if !netmask.Contains(nextIp) {
			break
		}
	}

	fmt.Println("=============================")
	fmt.Println(aliveIps)

}

func getNextIP(ip net.IP, inc uint) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
}
