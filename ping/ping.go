package ping

import (
	"net"
	"os"
	"time"

	"log/slog"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Status struct {
	Alive  bool
	IpAddr string
}

func Ping(targetIp net.IP) (Status, error) {

	status := Status{
		Alive:  false,
		IpAddr: targetIp.String(),
	}

	destAddr, err := net.ResolveIPAddr("ip", targetIp.String())

	if err != nil {
		slog.Error("error parsing ip: ", err)
		return status, err
	}

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")

	if err != nil {
		slog.Error("could not establish socker conn: ", err)
		return status, err
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
		return status, err
	}

	_, err = c.WriteTo(msgBytes, destAddr)
	if err != nil {
		slog.Error("could not write msg: ", err)
		return status, err
	}

	reply := make([]byte, 1500)

	err = c.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		slog.Error("error setting deadline: ", err)
		return status, err
	}

	n, addr, err := c.ReadFrom(reply)
	if err != nil {
		slog.Error("error receiving ICMP reply: ", err)
		return status, err
	}

	slog.Info("Received from: ", addr.String())
	replyMsg, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		slog.Error("error parsing ICMP reply: ", replyMsg)
		return status, err
	}

	switch replyMsg.Type {
	case ipv4.ICMPTypeEchoReply:
		slog.Info("Received ICMP Echo Reply from: ", destAddr.String())
		status.Alive = true
		return status, nil
	default:
		slog.Info("Received ICMP message of type: ", replyMsg.Type)
		return status, nil
	}
}
