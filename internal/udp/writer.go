package udp

import (
	"bufio"
	"chat2p2/internal/datagram"
	"chat2p2/internal/message"
	"chat2p2/internal/peers"
	"chat2p2/pkg/logger"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

func WriteLoop(conn *net.UDPConn) {
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

		dg := message.CreateMessage(input, datagram.TEXT)
		bytes, err := dg.Marshal()
		if err != nil {
			log.Error("Error converting message to bytes", zap.Error(err))
			continue
		}
		peersRepo := peers.NewStore()
		targets := peersRepo.GetAllAddresses() // Obtém os endereços de todos os usuários registrados
		Broadcast(conn, bytes, targets)        // Envia para todos na rede local
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

func Broadcast(conn *net.UDPConn, data []byte, targets []string) {
	log := logger.Get()
	var wg sync.WaitGroup
	errs := make(chan error, len(targets))

	for _, addrStr := range targets {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()

			addr, err := net.ResolveUDPAddr("udp", target)
			if err != nil {
				errs <- fmt.Errorf("resolve %s: %w", target, err)
				return
			}

			_, err = conn.WriteToUDP(data, addr)
			if err != nil {
				errs <- fmt.Errorf("write %s: %w", target, err)
			}
		}(addrStr)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		log.Error("Broadcast error", zap.Error(err))
	}
}
