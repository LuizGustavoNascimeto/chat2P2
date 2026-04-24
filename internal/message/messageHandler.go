package message

import (
	"chat2p2/internal/datagram"
	"chat2p2/pkg/logger"

	"go.uber.org/zap"
)

func MessageHandler(dg *datagram.Datagram, addr string) {
	log := logger.Get()
	log.Info("Received message", zap.String("message", dg.MessageText), zap.String("from", addr))

}

func CreateMessage(text string, t datagram.MessageType) *datagram.Datagram {
	if len(text) > 255 {
		text = text[:255] // Truncate message if it exceeds 255 bytes
	}

	nick := "User" // Placeholder for nickname, can be extended to include actual user management
	dg := &datagram.Datagram{
		MessageType: t,
		MessageSize: uint8(len(text)),
		MessageText: text,
		NickSize:    uint8(len(nick)),
		Nick:        nick,
	}

	log := logger.Get()
	log.Info("Created message datagram", zap.String("message", text))
	return dg
}
