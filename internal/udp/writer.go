package udp

import (
	"bufio"
	"chat2p2/internal/message"
	"chat2p2/internal/peers"
	"chat2p2/internal/sender"
	"chat2p2/pkg/logger"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"go.uber.org/zap"
)

func WriteLoop(conn *net.UDPConn, peersRepo *peers.Store) {
	log := logger.Get()
	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := ReadLine("> ", reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Info("stdin closed, exiting write loop")
				return // encerra em vez de spammar
			}
			log.Error("Error reading input", zap.Error(err))
			continue
		}

		dg := message.CreateMessage(input)
		bytes, err := dg.Marshal()
		if err != nil {
			log.Error("Error converting message to bytes", zap.Error(err))
			continue
		}
		targets := peersRepo.GetAllAddresses() // Obtém os endereços de todos os usuários registrados
		sender.Broadcast(conn, bytes, targets) // Envia para todos na rede local
	}

}

func ReadLine(prompt string, reader *bufio.Reader) (string, error) {
	if prompt != "" {
		fmt.Print(prompt)
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Error reading input: %w", err)
	}

	return strings.TrimRight(input, "\r\n"), nil // remove \n e \r\n (Windows)
}
