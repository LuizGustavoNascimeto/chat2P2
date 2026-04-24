package udp

import (
	"chat2p2/internal/datagram"
	"chat2p2/internal/message"
	"chat2p2/internal/users"
	"chat2p2/pkg/logger"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
)

// broadcastPing envia um ping UDP para todas as portas no intervalo [start, end].
func broadcastPing(conn *net.UDPConn, data []byte, start, end int) {
	log := logger.Get()
	var wg sync.WaitGroup
	for port := start; port <= end; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			target := fmt.Sprintf("127.0.0.1:%d", p)
			addr, err := net.ResolveUDPAddr("udp", target)
			if err != nil {
				log.Error("Failed to resolve UDP target", zap.Int("port", p), zap.Error(err))
				return
			}
			if _, err = conn.WriteToUDP(data, addr); err != nil {
				log.Error("Failed to send UDP ping", zap.Int("port", p), zap.Error(err))
			}
		}(port)
	}
	wg.Wait()
}

// collectEchoReplies lê respostas UDP durante o timeout e retorna os peers
// que responderam com ECHO.
func collectEchoReplies(conn *net.UDPConn, timeout time.Duration) []users.User {
	log := logger.Get()
	conn.SetReadDeadline(time.Now().Add(timeout))
	defer conn.SetReadDeadline(time.Time{})

	buf := make([]byte, 4096)
	var found []users.User

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			break
		}
		var response datagram.Datagram
		if err := json.Unmarshal(buf[:n], &response); err != nil {
			continue
		}
		if response.MessageType == datagram.ECHO {
			found = append(found, users.User{
				Name:     response.Nick,
				Addr:     remoteAddr.String(),
				LastSeen: time.Now().Unix(),
			})
			log.Info("Peer found", zap.String("addr", remoteAddr.String()))
		}
	}
	return found
}

// FindPeers descobre peers ativos nas portas 8000–9000 via ping/pong UDP.
func FindPeers(conn *net.UDPConn) []users.User {
	dg := message.CreateMessage("ping", datagram.ECHO)
	data, _ := json.Marshal(dg)

	go broadcastPing(conn, data, 8000, 9000)

	return collectEchoReplies(conn, 5*time.Second)
}
