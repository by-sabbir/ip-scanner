package ping

import (
	"net"
	"os"
	"time"

	"golang.org/x/exp/slog"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func ping() (bool, error) {
	targetIp := "192.168.0.1"
	destAddr, err := net.ResolveIPAddr("ip", targetIp)

	if err != nil {
		slog.Error("error parsing ip: ", err)
		return false, err
	}

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")

	if err != nil {
		slog.Error("could not establish socker conn: ", err)
		return false, err
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
		return false, err
	}

	_, err = c.WriteTo(msgBytes, destAddr)
	if err != nil {
		slog.Error("could not write msg: ", err)
		return false, err
	}

	reply := make([]byte, 1500)

	err = c.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		slog.Error("error setting deadline: ", err)
		return false, err
	}

	n, addr, err := c.ReadFrom(reply)
	if err != nil {
		slog.Error("error receiving ICMP reply: ", err)
		return false, err
	}

	slog.Info("Received from: ", addr.String())
	replyMsg, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		slog.Error("error parsing ICMP reply: ", replyMsg)
		return false, err
	}

	switch replyMsg.Type {
	case ipv4.ICMPTypeEchoReply:
		slog.Info("Received ICMP Echo Reply from: ", destAddr.String())
		return true, nil
	default:
		slog.Info("Received ICMP message of type: ", replyMsg.Type)
		return true, nil
	}
}
