Projeto: Chat P2P UDP
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25

1. Como compilar
- Requisito: Go 1.25+ instalado.
- No diretório do projeto, execute:
  go mod tidy
  go build -o chat2p2 ./cmd

2. Como executar
- Em um terminal:
  ./chat2p2 <apelido>
- Em Windows PowerShell/cmd:
  .\chat2p2.exe <apelido>
- Para simular múltiplos peers locais, abra 2+ terminais e execute o binário em cada um com apelidos diferentes.

3. Bibliotecas usadas
- Padrão da linguagem Go (standard library):
  fmt, net, os, strings, time, sync, bufio, errors, io, strconv
- Não padrão:
  go.uber.org/zap (v1.27.1): logging estruturado de alta performance.
  go.uber.org/multierr (indireta): agregação de múltiplos erros (dependência do zap).

4. Exemplo de uso
- Terminal 1:
  ./chat2p2 Alice
- Terminal 2:
  ./chat2p2 Bob
- Comandos de mensagem:
  Olá pessoal
  /emoji :D
  /url https://golang.org
