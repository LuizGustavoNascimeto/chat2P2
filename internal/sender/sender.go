package sender

import (
	"chat2p2/pkg/logger"
	"fmt"
	"net"
	"sync"

	"go.uber.org/zap"
)

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
