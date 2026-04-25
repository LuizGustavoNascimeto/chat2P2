package message

import (
	"chat2p2/internal/datagram"
	"chat2p2/pkg/logger"
	"os"
	"strings"

	"go.uber.org/zap"
)

func CreateMessage(text string) *datagram.Datagram {
	if len(text) > 255 {
		text = text[:255] // Truncate message if it exceeds 255 bytes
	}

	args := os.Args
	var nick string
	if len(args) > 1 {
		nick = args[1]
	} else {
		nick = "???"
	}

	t, content := processMessage(text)

	dg := &datagram.Datagram{
		Type:        t,
		MessageSize: uint8(len(content)),
		MessageText: content,
		NickSize:    uint8(len(nick)),
		Nick:        nick,
	}

	log := logger.Get()
	log.Debug("Created message datagram", zap.String("message", text))
	return dg
}

// cria um echo genérica da vida
func CreateEcho(text string) *datagram.Datagram {
	args := os.Args
	var nick string
	if len(args) > 1 {
		nick = args[1]
	} else {
		nick = "???"
	}
	dg := &datagram.Datagram{
		Type:        datagram.ECHO,
		MessageSize: uint8(len(text)),
		MessageText: text,
		NickSize:    uint8(len(nick)),
		Nick:        nick,
	}

	return dg
}

// ele vai ver se a mensagem é um comando como /url texto ou /emoji texto para definir o type
// emoji vai ser tratado como texto por hora
func processMessage(text string) (datagram.MessageType, string) {
	if text[0] == '/' && len(strings.Fields(text)) > 1 {
		return processCommand(text)
	}
	return datagram.TEXT, text

}

func processCommand(text string) (datagram.MessageType, string) {
	fields := strings.Fields(text)
	switch fields[0] {
	case "/url":
		return datagram.URL, fields[1]
	case "/emoji":
		return datagram.EMOJI, fields[1]
	}
	return datagram.TEXT, text
}
