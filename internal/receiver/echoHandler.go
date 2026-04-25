/*
Descrição: Trata datagramas ECHO para descoberta e confirmação de peers ativos.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package receiver

import (
	"chat2p2/internal/env"
	"chat2p2/internal/peers"
	"chat2p2/internal/sender"
	"chat2p2/pkg/utils"
	"net"
)

// EchoHandler processa pacote ECHO, registra peer novo e responde quando necessário.
// Entradas:  endereço remoto, repositório de peers e conexão UDP.
// Saída: erro em falhas de parsing/registro/envio; nil em sucesso.
func EchoHandler(addr string, peersRepo *peers.Store, conn *net.UDPConn) error {
	myPort := env.GetPort()
	addrPort, err := utils.ExtractPort(addr)
	if err != nil {
		return err
	}
	if myPort == addrPort {
		// Ignorar mensagens ECHO enviadas por mim mesmo
		return nil
	}
	peer, err := peersRepo.Read(addr)
	//se o peer não for encontrado é uma descoberta de peer, então registra ele
	if err == peers.ErrPeerNotFound {
		peer, err = peersRepo.Create(addr)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if peer.Waiting {
		peer.Waiting = false
	} else {
		sender.ResponseEcho(addr, conn)
	}
	return nil
}
