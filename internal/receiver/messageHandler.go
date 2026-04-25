package receiver

import (
	"chat2p2/internal/datagram"
	"chat2p2/internal/peers"
	"chat2p2/pkg/logger"
	"net"
	"time"

	"go.uber.org/zap"
)

func MessageHandler(dg *datagram.Datagram, addr string, conn *net.UDPConn, peersRepo *peers.Store) {
	log := logger.Get()
	log.Debug(dg.String(), zap.String("from", addr))
	peersRepo.UpdateLastSeen(addr, time.Now().Unix())

	switch dg.Type {
	case datagram.TEXT:
		TextHandler(dg, addr)
	case datagram.ECHO:
		EchoHandler(dg, addr, peersRepo, conn)
	case datagram.EMOJI:
		EmojiHandler(dg, addr)
	case datagram.URL:
		URLHandler(dg, addr)

	default:
		log.Warn("Received datagram with unknown type", zap.Int("type", int(dg.Type)))
	}

}
