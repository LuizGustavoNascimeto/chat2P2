/*
Descrição: Define formato de datagrama UDP e faz serialização/desserialização de mensagens.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package datagram

import (
	"chat2p2/pkg/logger"

	"go.uber.org/zap"
)

type MessageType uint8

const (
	TEXT  MessageType = 1
	EMOJI MessageType = 2
	URL   MessageType = 3
	ECHO  MessageType = 4
)

type Datagram struct {
	Type        MessageType
	NickSize    uint8
	Nick        string
	MessageSize uint8
	MessageText string
}

// Unmarshal converte bytes recebidos em estrutura Datagram validada.
// Entradas: slice de bytes no formato [tipo|nickSize|nick|messageSize|message].
// Saída: ponteiro para Datagram e erro em caso de payload inválido.
func Unmarshal(data []byte) (*Datagram, error) {
	// Representa posição atual durante parsing do payload.
	log := logger.Get()
	var cursor int

	if len(data) < 4 {
		return nil, logger.LogError("INVALID_DATAGRAM: datagram must be at least 3 bytes long")
	}
	// Lê tipo e tamanho do nickname.
	dg := &Datagram{
		Type:     MessageType(data[0]),
		NickSize: data[1],
	}
	cursor += 2

	// Lê nickname.
	if dg.NickSize > 0 && dg.NickSize <= 64 {
		dg.Nick = string(data[cursor : cursor+int(dg.NickSize)])
		cursor += int(dg.NickSize)
	} else {
		return nil, logger.LogError("INVALID_DATAGRAM: invalid nickname size")
	}

	// Lê tamanho e conteúdo da mensagem.
	dg.MessageSize = data[cursor]
	cursor++
	dg.MessageText = string(data[cursor : cursor+int(dg.MessageSize)])
	log.Debug("Parsed datagram", zap.Uint8("type", uint8(dg.Type)), zap.String("nick", dg.Nick), zap.String("message", dg.MessageText))

	return dg, nil

}

// Marshal converte Datagram em bytes prontos para envio em UDP.
// Entradas: receptor Datagram com campos Type, Nick e MessageText preenchidos.
// Saída: slice de bytes serializado e erro se houver inconsistência de tamanho.
func (dg *Datagram) Marshal() ([]byte, error) {
	if dg.NickSize > 64 {
		return nil, logger.LogError("INVALID_DATAGRAM: nickname size exceeds 64 bytes")
	}
	//messageType (1 byte) + nickSize (1 byte) + nick (nickSize bytes) + messageSize (1 byte) + messageText (messageSize bytes)
	data := make([]byte, 1+1+dg.NickSize+1+dg.MessageSize)
	data[0] = byte(dg.Type)
	data[1] = byte(dg.NickSize)
	copy(data[2:2+dg.NickSize], []byte(dg.Nick))
	data[2+dg.NickSize] = byte(dg.MessageSize)
	copy(data[3+dg.NickSize:], []byte(dg.MessageText))
	return data, nil

}

// TypeToString converte MessageType para representação textual legível.
// Entradas: tipo da mensagem.
// Saída: string de identificação do tipo.
func TypeToString(messageType MessageType) string {
	switch messageType {
	case TEXT:
		return "TEXT"
	case EMOJI:
		return "EMOJI"
	case URL:
		return "URL"
	case ECHO:
		return "ECHO"
	default:
		return "UNKNOWN"
	}
}

// String retorna representação simplificada do datagrama para logs.
// Entradas: receptor Datagram.
// Saída: string no formato TYPE/NICK/MESSAGE.
func (dg *Datagram) String() string {
	return TypeToString(dg.Type) + "/" + dg.Nick + "/" + string(dg.MessageText)
}
