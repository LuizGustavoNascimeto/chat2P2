package chat2p2

import (
	"chat2p2/internal/udp"
	"chat2p2/internal/users"
	"chat2p2/pkg/logger"
	"chat2p2/pkg/utils"
	"fmt"
	"net"

	"go.uber.org/zap"
)

func main() {
	logger.Init("development")
	defer logger.Sync()

	log := logger.Get()

	port, err := utils.FindAvailablePort(8000, 9000)
	if err != nil {
		log.Fatal("Failed to find an available port", zap.Error(err))
	}
	log.Info("Starting UDP chat", zap.Int("port", port))

	conn, err := udp.StartUDP(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Failed to start UDP server", zap.Error(err), zap.Int("port", port))
	}
	defer conn.Close()

	findUsers(conn)

	// goroutine de leitura
	go udp.ReadLoop(conn)

	udp.WriteLoop(conn)

}

func findUsers(conn *net.UDPConn) {
	log := logger.Get()
	peers := udp.FindPeers(conn)
	if len(peers) == 0 {
		log.Info("No peers found")
	} else {
		log.Info("Peers found", zap.Int("count", len(peers)))
		usersRepo := users.NewStore()
		for _, peer := range peers {
			if err := usersRepo.Create(peer); err != nil {
				log.Error("Failed to create user", zap.String("addr", peer.Addr), zap.Error(err))
			} else {
				log.Info("Peer registered", zap.String("name", peer.Name), zap.String("addr", peer.Addr))
			}
		}
	}
}
