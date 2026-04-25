/*
Descrição: Renderiza mensagens do tipo URL em formato de hyperlink para terminal compatível.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package receiver

import (
	"chat2p2/internal/datagram"
	"fmt"
)

// URLHandler imprime URL recebida com escape sequence de hyperlink ANSI.
// Entradas: datagrama com URL e endereço remoto do emissor.
// Saída: nenhuma; efeito colateral de escrita em stdout.
func URLHandler(receivedDatagram *datagram.Datagram, addr string) {
	_ = addr
	urlText := receivedDatagram.MessageText
	fmt.Printf("%s: %s\n", receivedDatagram.Nick, hyperlink(urlText, urlText))

}

// hyperlink monta sequência ANSI OSC 8 para link clicável no terminal.
// Entradas: URL de destino e texto de exibição.
// Saída: string com sequência ANSI formatada.
func hyperlink(url, text string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, text)
}

// func main() {
// 	url := "https://golang.org"

// 	// Link simples
// 	fmt.Println(hyperlink(url, url))

// 	// Link com texto customizado
// 	fmt.Println(hyperlink(url, "Clique aqui para acessar o Go!"))
// }
