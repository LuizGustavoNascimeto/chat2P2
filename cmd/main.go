/*
Descrição: Inicializa chat P2P via UDP, descobre peers e inicia loops de leitura e escrita.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package main

import (
	"chat2p2/internal/env"
	"chat2p2/internal/peers"
	"chat2p2/internal/udp"
	"chat2p2/pkg/logger"
	"net"

	"go.uber.org/zap"
)

// main configura logger, abre socket UDP, descobre peers e inicia loops de I/O.
// Entradas: argumentos de linha de comando opcionais (apelido) lidos por camadas internas.
// Saída: nenhuma; encerra processo em falha crítica de inicialização.
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

	// Descobre peers na rede local e registra no repositório.
	discoverAndRegisterPeers(conn, peersRepo)

	// Inicia goroutine de leitura.
	go udp.ReadLoop(conn, peersRepo)

	// Mantém loop principal de escrita.
	udp.WriteLoop(conn, peersRepo)

}

// discoverAndRegisterPeers encontra peers ativos e registra cada endereço no repositório.
// Entradas: conexão UDP ativa e repositório de peers.
// Saída: nenhuma; efeitos colaterais no repositório e logs.
func discoverAndRegisterPeers(conn *net.UDPConn, peersRepo *peers.Store) {
	log := logger.Get()
	foundPeers := udp.FindPeers(conn)
	if len(foundPeers) == 0 {
		log.Info("No peers found")
	} else {
		log.Info("Peers found", zap.Int("count", len(foundPeers)))
		for _, peer := range foundPeers {
			if _, err := peersRepo.Create(peer.Addr); err != nil {
				log.Error("Failed to create user", zap.String("addr", peer.Addr), zap.Error(err))
			}
		}
	}
}
