package sender

import (
	"chat2p2/internal/message"
	"net"
)

func ResponseEcho(addr string, conn *net.UDPConn) error {
	dg := message.CreateEcho("pong")
	data, err := dg.Marshal()
	if err != nil {
		return err
	}
	fullAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	_, err = conn.WriteToUDP(data, fullAddr)
	if err != nil {
		return err
	}

	return nil
}
