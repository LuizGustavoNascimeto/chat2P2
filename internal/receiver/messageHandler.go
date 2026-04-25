/*
Descrição: Roteia datagramas recebidos para handlers por tipo de mensagem.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package receiver

import (
	"chat2p2/internal/datagram"
	"chat2p2/internal/peers"
	"chat2p2/pkg/logger"
	"net"
	"time"

	"go.uber.org/zap"
)

// MessageHandler atualiza presença do peer e delega processamento por tipo.
// Entradas: datagrama recebido, endereço remoto, conexão UDP e repositório de peers.
// Saída: nenhuma; efeitos colaterais em logs, store e stdout.
func MessageHandler(receivedDatagram *datagram.Datagram, addr string, conn *net.UDPConn, peersRepo *peers.Store) {
	log := logger.Get()
	log.Debug(receivedDatagram.String(), zap.String("from", addr))
	peersRepo.UpdateLastSeen(addr, time.Now().Unix())

	switch receivedDatagram.Type {
	case datagram.TEXT:
		TextHandler(receivedDatagram, addr)
	case datagram.ECHO:
		EchoHandler(addr, peersRepo, conn)
	case datagram.EMOJI:
		EmojiHandler(receivedDatagram, addr)
	case datagram.URL:
		URLHandler(receivedDatagram, addr)

	default:
		log.Warn("Received datagram with unknown type", zap.Int("type", int(receivedDatagram.Type)))
	}

}
