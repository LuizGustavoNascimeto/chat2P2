package datagram

import (
	"chat2p2/pkg/logger"
)

type MessageType int8

const (
	TEXT  MessageType = 1
	EMOJI MessageType = 2
	URL   MessageType = 3
	ECHO  MessageType = 4
)

type Datagram struct {
	MessageType MessageType
	NickSize    uint8
	Nick        string
	MessageSize uint8
	MessageText string
}

func Unmarshal(data []byte) (*Datagram, error) {
	//represents the current position in the datagram
	var pos int = 0

	if len(data) < 4 {
		return nil, logger.LogError("INVALID_DATAGRAM: datagram must be at least 3 bytes long")
	}
	//set o nicksize e o messageType
	dg := &Datagram{
		MessageType: MessageType(data[0]),
		NickSize:    data[1],
	}
	pos += 2

	//set o nick
	if dg.NickSize > 0 && dg.NickSize <= 64 {
		dg.Nick = string(data[pos : pos+int(dg.NickSize)])
		pos += int(dg.NickSize)
	} else {
		return nil, logger.LogError("INVALID_DATAGRAM: invalid nickname size")
	}

	//set o messagesize e o messagetext
	dg.MessageSize = data[pos]
	pos++
	dg.MessageText = string(data[pos : pos+int(dg.MessageSize)])

	return dg, nil

}
func (dg *Datagram) Marshal() ([]byte, error) {
	if dg.NickSize > 64 {
		return nil, logger.LogError("INVALID_DATAGRAM: nickname size exceeds 64 bytes")
	}
	//messageType (1 byte) + nickSize (1 byte) + nick (nickSize bytes) + messageSize (1 byte) + messageText (messageSize bytes)
	data := make([]byte, 1+1+dg.NickSize+1+dg.MessageSize)
	data[0] = byte(dg.MessageType)
	data[1] = byte(dg.NickSize)
	copy(data[2:2+dg.NickSize], []byte(dg.Nick))
	data[2+dg.NickSize] = byte(dg.MessageSize)
	copy(data[3+dg.NickSize:], []byte(dg.MessageText))
	return data, nil

}
