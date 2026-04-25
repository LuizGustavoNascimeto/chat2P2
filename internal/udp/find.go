/*
Descrição: Descobre peers ativos por broadcast de ping/pong em intervalo de portas UDP.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

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

// broadcastPing envia ping UDP para todas portas no intervalo informado.
// Entradas: conexão UDP, payload serializado, porta inicial e porta final.
// Saída: nenhuma; erros de envio são registrados em log.
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

// collectEchoReplies lê respostas ECHO por janela de tempo e monta lista de peers.
// Entradas: conexão UDP e tempo máximo de coleta.
// Saída: slice com peers que responderam com mensagem ECHO.
func collectEchoReplies(conn *net.UDPConn, timeout time.Duration) []peers.Peer {
	conn.SetReadDeadline(time.Now().Add(timeout))
	defer conn.SetReadDeadline(time.Time{})

	buf := make([]byte, 4096)
	var found []peers.Peer

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			break
		}
		responseDatagram, err := datagram.Unmarshal(buf[:n])
		if err != nil {
			continue
		}

		if responseDatagram.Type == datagram.ECHO {
			found = append(found, peers.Peer{
				Addr:     remoteAddr.String(),
				LastSeen: time.Now().Unix(),
			})
		}
	}
	return found
}

// FindPeers executa ciclo ping/pong para descobrir peers ativos.
// Entradas: conexão UDP ativa.
// Saída: slice de peers encontrados; nil em falhas de serialização.
func FindPeers(conn *net.UDPConn) []peers.Peer {
	pingDatagram := message.CreateEcho("ping")
	data, err := pingDatagram.Marshal()
	if err != nil {
		logger.Get().Error("Failed to marshal ping datagram", zap.Error(err))
		return nil
	}

	go broadcastPing(conn, data, 8000, 9000)

	return collectEchoReplies(conn, 5*time.Second)
}
