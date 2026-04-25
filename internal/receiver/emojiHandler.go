/*
Descrição: Converte payload de emoji em símbolo Unicode e imprime no terminal.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package receiver

import (
	"chat2p2/internal/datagram"
	"fmt"
	"strings"
)

// EmojiHandler interpreta texto como atalho de emoji e exibe saída formatada.
// Entradas: datagrama de mensagem e endereço remoto do emissor.
// Saída: nenhuma; efeito colateral de escrita em stdout.
func EmojiHandler(receivedDatagram *datagram.Datagram, addr string) {
	_ = addr
	text := strings.TrimSpace(receivedDatagram.MessageText)
	emojiMap := map[string]string{
		":D":  "😄",
		":)":  "🙂",
		":(":  "🙁",
		";)":  "😉",
		":P":  "😛",
		":O":  "😮",
		":'(": "😢",
		"<3":  "❤️",
		":|":  "😐",
		":*":  "😘",
	}

	emoji, ok := emojiMap[text]
	if !ok {
		emoji = "❓"
	}

	fmt.Printf("%s: %s\n", receivedDatagram.Nick, emoji)
}
