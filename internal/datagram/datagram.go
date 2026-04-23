package datagram

import "fmt"

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

func ParseDatagram(data []byte) (*Datagram, error) {
	//represents the current position in the datagram
	var pos int = 0

	if len(data) < 4 {
		return nil, fmt.Errorf("INVALID_DATAGRAM: datagram must be at least 3 bytes long")
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
		return nil, fmt.Errorf("INVALID_DATAGRAM: invalid nickname size")
	}

	//set o messagesize e o messagetext
	dg.MessageSize = data[pos]
	pos++
	dg.MessageText = string(data[pos : pos+int(dg.MessageSize)])

	return dg, nil

}
