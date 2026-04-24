package utils

import (
	"chat2p2/pkg/logger"
	"fmt"
	"net"

	"go.uber.org/zap"
)

func FindAvailablePort(start, end int) (int, error) {
	if start < 1 || end > 65535 || start > end {
		return 0, logger.LogError(fmt.Sprintf("INVALID_PORT_RANGE: intervalo de portas inválido: %d-%d", start, end))
	}

	for port := start; port <= end; port++ {
		addr := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			continue // porta ocupada ou sem permissão, tenta a próxima
		}
		ln.Close()
		logger.Get().Info("Available port found", zap.Int("port", port), zap.Int("start", start), zap.Int("end", end))
		return port, nil
	}

	return 0, logger.LogError(fmt.Sprintf("NO_AVAILABLE_PORT: nenhuma porta disponível no intervalo %d-%d", start, end))
}
