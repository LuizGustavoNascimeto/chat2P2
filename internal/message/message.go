/*
Descrição: Constrói datagramas de chat e interpreta comandos de entrada do usuário.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package message

import (
	"chat2p2/internal/datagram"
	"chat2p2/pkg/logger"
	"os"
	"strings"

	"go.uber.org/zap"
)

// CreateMessage cria datagrama padrão a partir de texto digitado pelo usuário.
// Entradas: texto bruto do terminal.
// Saída: ponteiro para Datagram pronto para serialização e envio.
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

	messageType, content := processMessage(text)

	dg := &datagram.Datagram{
		Type:        messageType,
		MessageSize: uint8(len(content)),
		MessageText: content,
		NickSize:    uint8(len(nick)),
		Nick:        nick,
	}

	log := logger.Get()
	log.Debug("Created message datagram", zap.String("message", text))
	return dg
}

// CreateEcho cria datagrama ECHO para descoberta e resposta de peers.
// Entradas: texto do payload de echo (ex.: ping/pong).
// Saída: ponteiro para Datagram do tipo ECHO.
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

// processMessage identifica se entrada é comando e resolve tipo de mensagem.
// Entradas: texto completo digitado pelo usuário.
// Saída: tipo de mensagem e conteúdo final a ser enviado.
func processMessage(text string) (datagram.MessageType, string) {
	if len(text) > 0 && text[0] == '/' && len(strings.Fields(text)) > 1 {
		return processCommand(text)
	}
	return datagram.TEXT, text

}

// processCommand converte comandos de barra em tipo e payload compatíveis.
// Entradas: texto de comando (ex.: "/url https://site").
// Saída: tipo de mensagem correspondente e conteúdo do comando.
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
