package datagram

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
	Message     string
}
