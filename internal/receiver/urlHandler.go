package receiver

import (
	"chat2p2/internal/datagram"
	"fmt"
)

func URLHandler(dg *datagram.Datagram, addr string) {
	text := dg.MessageText
	fmt.Printf("%s: %s\n", dg.Nick, hyperlink(text, text))

}

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
