package receiver

import (
	"chat2p2/internal/datagram"
	"chat2p2/internal/peers"
	"chat2p2/pkg/logger"

	"go.uber.org/zap"
)

func MessageHandler(dg *datagram.Datagram, addr string) {
	peersRepo := peers.NewStore()
	log := logger.Get()
	log.Info(dg.String(), zap.String("from", addr))

	switch dg.Type {
	case datagram.TEXT:

	case datagram.ECHO:
		EchoHandler(dg, addr, peersRepo)
	default:
		log.Warn("Received datagram with unknown type", zap.Int("type", int(dg.Type)))
	}

}
