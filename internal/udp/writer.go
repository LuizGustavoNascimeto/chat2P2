package udp

import (
	"bufio"
	"chat2p2/internal/datagram"
	"chat2p2/internal/message"
	"chat2p2/internal/users"
	"chat2p2/pkg/logger"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

func WriteLoop(conn *net.UDPConn) {
	log := logger.Get()

	for {
		input, err := ReadLine("> ")
		if err != nil {
			log.Error("Error reading input", zap.Error(err))
			continue
		}

		dg := message.CreateMessage(input, datagram.TEXT)
		bytes, err := dg.Marshal()
		if err != nil {
			log.Error("Error converting message to bytes", zap.Error(err))
			continue
		}
		usersRepo := users.NewStore()
		targets := usersRepo.GetAllAddresses() // Obtém os endereços de todos os usuários registrados
		Broadcast(conn, bytes, targets)        // Envia para todos na rede local
	}

}

func ReadLine(prompt string) (string, error) {
	if prompt != "" {
		fmt.Print(prompt)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("erro ao ler entrada: %w", err)
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
