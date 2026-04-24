package env

import (
	"chat2p2/pkg/utils"
	"fmt"
	"sync"
)

var (
	addr string
	once sync.Once
)

func GetPort() string {
	once.Do(func() {
		port, err := utils.FindAvailablePortUDP(8000, 9000)
		if err != nil {
			panic(err)
		}
		addr = fmt.Sprintf(":%d", port)
	})
	return addr
}
