package main

import (
	"net"
	"os"
	"time"

	"log/slog"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	targetIp := "192.168.0.1"
	destAddr, err := net.ResolveIPAddr("ip", targetIp)

	if err != nil {
		slog.Error("error parsing ip: ", err)
	}

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")

	if err != nil {
		slog.Error("could not establish socker conn: ", err)
	}
	defer c.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid(),
			Seq:  1,
			Data: []byte("ping"),
		},
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		slog.Error("error marshalling: ", err)
	}

	_, err = c.WriteTo(msgBytes, destAddr)
	if err != nil {
		slog.Error("could not write msg: ", err)
	}

	reply := make([]byte, 1500)

	err = c.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		slog.Error("error setting deadline: ", err)
	}

	n, addr, err := c.ReadFrom(reply)
	if err != nil {
		slog.Error("error receiving ICMP reply: ", err)
	}

	slog.Info("Received from: ", addr.String())
	replyMsg, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		slog.Error("error parsing ICMP reply: ", replyMsg)
	}

	switch replyMsg.Type {
	case ipv4.ICMPTypeEchoReply:
		slog.Info("Received ICMP Echo Reply from %s\n", destAddr.String())
	default:
		slog.Info("Received ICMP message of type %v\n", replyMsg.Type)
	}

}
