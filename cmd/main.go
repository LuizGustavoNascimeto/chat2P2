package main

import (
	"chat2p2/internal/env"
	"chat2p2/internal/peers"
	"chat2p2/internal/udp"
	"chat2p2/pkg/logger"
	"net"

	"go.uber.org/zap"
)

func main() {
	logger.Init("development")
	defer logger.Sync()

	log := logger.Get()

	port := env.GetPort()

	log.Info("Starting UDP chat", zap.String("port", port))

	conn, err := udp.StartUDP(port)
	if err != nil {
		log.Fatal("Failed to start UDP server", zap.Error(err), zap.String("port", port))
	}
	defer conn.Close()
	peersRepo := peers.NewStore()

	// Descobrir peers na rede local e registrá-los no repositório de usuários
	setPeers(conn, peersRepo)

	// goroutine de leitura
	go udp.ReadLoop(conn, peersRepo)

	// loop de escrita
	udp.WriteLoop(conn, peersRepo)

}

func setPeers(conn *net.UDPConn, peersRepo *peers.Store) {
	log := logger.Get()
	p := udp.FindPeers(conn)
	if len(p) == 0 {
		log.Info("No peers found")
	} else {
		log.Info("Peers found", zap.Int("count", len(p)))
		for _, peer := range p {
			if _, err := peersRepo.Create(peer.Addr); err != nil {
				log.Error("Failed to create user", zap.String("addr", peer.Addr), zap.Error(err))
			}
		}
	}
}
