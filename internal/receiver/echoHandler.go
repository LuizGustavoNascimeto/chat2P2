package receiver

import (
	"chat2p2/internal/datagram"
	"chat2p2/internal/env"
	"chat2p2/internal/peers"
	"chat2p2/internal/sender"
	"chat2p2/pkg/utils"
	"net"
)

func EchoHandler(dg *datagram.Datagram, addr string, peersRepo *peers.Store, conn *net.UDPConn) error {
	myPort := env.GetPort()
	addrPort, err := utils.ExtractPort(addr)
	if err != nil {
		return err
	}
	if myPort == addrPort {
		// Ignorar mensagens ECHO enviadas por mim mesmo
		return nil
	}
	peer, err := peersRepo.Read(addr)
	//se o peer não for encontrado é uma descoberta de peer, então registra ele
	if err == peers.ErrPeerNotFound {
		peer, err = peersRepo.Create(addr)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if peer.Waiting {
		peer.Waiting = false
	} else {
		sender.ResponseEcho(addr, conn)
	}
	return nil
}
