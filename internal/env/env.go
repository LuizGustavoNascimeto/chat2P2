/*
Descrição: Centraliza resolução lazy de porta local usada pelo processo de chat UDP.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

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

// GetPort retorna endereço local no formato ":porta", inicializado apenas uma vez.
// Entradas: nenhuma.
// Saída: string contendo porta escolhida no range configurado.
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
