/*
Descrição: Mantém loop de leitura UDP, faz parsing de datagramas e delega para handlers.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package udp

import (
	"chat2p2/internal/datagram"
	"chat2p2/internal/peers"
	"chat2p2/internal/receiver"
	"chat2p2/pkg/logger"
	"net"

	"go.uber.org/zap"
)

const MaxPacketSize = 322

// ReadLoop escuta socket UDP continuamente e despacha mensagens recebidas.
// Entradas: conexão UDP ativa e repositório de peers.
// Saída: nenhuma; loop contínuo até encerramento da conexão.
func ReadLoop(conn *net.UDPConn, peersRepo *peers.Store) {
	defer conn.Close()
	log := logger.Get()

	buf := make([]byte, MaxPacketSize) // cada leitura, um pacote completo
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Error("Error reading from UDP", zap.Error(err))
			continue
		}

		// processa o pacote recebido (buf[:n]) e remoteAddr
		dg, err := datagram.Unmarshal(buf[:n])
		if err != nil {
			log.Error("Error parsing datagram", zap.Error(err))
			continue
		}
		receiver.MessageHandler(dg, remoteAddr.String(), conn, peersRepo)
	}
}
