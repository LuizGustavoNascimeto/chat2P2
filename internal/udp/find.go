package udp

import (
	"chat2p2/internal/datagram"
	"chat2p2/internal/message"
	"chat2p2/internal/peers"
	"chat2p2/pkg/logger"
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
func collectEchoReplies(conn *net.UDPConn, timeout time.Duration) []peers.Peer {
	log := logger.Get()
	conn.SetReadDeadline(time.Now().Add(timeout))
	defer conn.SetReadDeadline(time.Time{})

	buf := make([]byte, 4096)
	var found []peers.Peer

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			break
		}
		response, err := datagram.Unmarshal(buf[:n])
		if err != nil {
			continue
		}
		if response.Type == datagram.ECHO {
			found = append(found, peers.Peer{
				Addr:     remoteAddr.String(),
				LastSeen: time.Now().Unix(),
			})
			log.Info("Peer found", zap.String("addr", remoteAddr.String()))
		}
	}
	return found
}

// FindPeers descobre peers ativos nas portas 8000–9000 via ping/pong UDP.
func FindPeers(conn *net.UDPConn) []peers.Peer {
	dg := message.CreateMessage("ping", datagram.ECHO)
	data, err := dg.Marshal()
	if err != nil {
		logger.Get().Error("Failed to marshal ping datagram", zap.Error(err))
		return nil
	}

	go broadcastPing(conn, data, 8000, 9000)

	return collectEchoReplies(conn, 5*time.Second)
}
