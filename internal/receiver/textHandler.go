/*
Descrição: Exibe mensagens de texto simples recebidas via UDP no terminal.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package receiver

import (
	"chat2p2/internal/datagram"
	"fmt"
)

// TextHandler imprime mensagem textual de peer remoto.
// Entradas: datagrama de texto e endereço remoto do emissor.
// Saída: nenhuma; efeito colateral de escrita em stdout.
func TextHandler(receivedDatagram *datagram.Datagram, addr string) {
	_ = addr
	fmt.Printf("%s: %s\n", receivedDatagram.Nick, receivedDatagram.MessageText)
}
