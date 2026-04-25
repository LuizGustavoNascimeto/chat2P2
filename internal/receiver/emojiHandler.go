package receiver

import (
	"chat2p2/internal/datagram"
	"fmt"
	"strings"
)

func EmojiHandler(dg *datagram.Datagram, addr string) {
	text := strings.TrimSpace(dg.MessageText)
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

	fmt.Printf("%s: %s\n", dg.Nick, emoji)
}
