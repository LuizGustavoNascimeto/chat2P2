/*
Descrição: Realiza envio concorrente de datagramas UDP para múltiplos peers.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package sender

import (
	"chat2p2/pkg/logger"
	"fmt"
	"net"
	"sync"

	"go.uber.org/zap"
)

// Broadcast envia payload para todos os alvos informados em paralelo.
// Entradas: conexão UDP, bytes serializados e lista de endereços destino.
// Saída: nenhuma; erros são registrados em log.
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
