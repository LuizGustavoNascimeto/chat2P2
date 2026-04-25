/*
Descrição: Envia resposta ECHO (pong) para endereço remoto durante descoberta de peers.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package sender

import (
	"chat2p2/internal/message"
	"net"
)

// ResponseEcho envia datagrama ECHO de resposta para peer solicitante.
// Entradas: endereço remoto de destino e conexão UDP ativa.
// Saída: erro quando serialização, resolução de endereço ou envio falhar.
func ResponseEcho(addr string, conn *net.UDPConn) error {
	echoDatagram := message.CreateEcho("pong")
	data, err := echoDatagram.Marshal()
	if err != nil {
		return err
	}
	fullAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	_, err = conn.WriteToUDP(data, fullAddr)
	if err != nil {
		return err
	}

	return nil
}
