/*
Descrição: Cria e inicializa conexão UDP de escuta do processo local.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package udp

import "net"

// StartUDP abre socket UDP vinculado ao endereço de escuta informado.
// Entradas: endereço de escuta no formato host:porta ou :porta.
// Saída: conexão UDP ativa e erro em falhas de resolução/bind.
func StartUDP(listenAddr string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
