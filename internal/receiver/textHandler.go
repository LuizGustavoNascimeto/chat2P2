package receiver

import (
	"chat2p2/internal/datagram"
	"fmt"
)

// se for texto é só printar :) e atualizar o lastSeen do peer
func TextHandler(dg *datagram.Datagram, addr string) {
	fmt.Printf("%s: %s\n", dg.Nick, dg.MessageText)
}
