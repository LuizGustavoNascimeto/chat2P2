/*
Descrição: Reúne utilitários de rede para descobrir portas UDP livres e extrair porta de endereço.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package utils

import (
	"bufio"
	"chat2p2/pkg/logger"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// FindAvailablePortUDP retorna primeira porta UDP livre no intervalo [start, end].
// Entradas: limite inicial e final de porta para busca.
// Saída: número da porta disponível ou erro quando range for inválido/sem portas livres.
// No Linux tenta ler /proc/net/udp para resolver em O(1) syscalls.
// Em outros sistemas faz fallback iterativo com ListenPacket.
func FindAvailablePortUDP(start, end int) (int, error) {
	if start < 1 || end > 65535 || start > end {
		return 0, logger.LogError(fmt.Sprintf(
			"INVALID_PORT_RANGE: intervalo inválido: %d-%d", start, end,
		))
	}

	if port, ok := findViaProc(start, end); ok {
		return port, nil
	}

	return findViaListen(start, end)
}

// findViaProc lê /proc/net/udp e /proc/net/udp6 para detectar portas ocupadas.
// Entradas: limite inicial e final de porta para busca.
// Saída: porta disponível e flag booleana indicando sucesso na estratégia /proc.
// Muito mais rápido que tentar bind em cada porta.
func findViaProc(start, end int) (int, bool) {
	f, err := os.Open("/proc/net/udp")
	if err != nil {
		return 0, false // não-Linux, usa fallback
	}
	defer f.Close()

	occupied := make(map[int]struct{})
	scanner := bufio.NewScanner(f)
	scanner.Scan() // pula o cabeçalho

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		// campo local_address: "0.0.0.0:hex_port"
		parts := strings.SplitN(fields[1], ":", 2)
		if len(parts) != 2 {
			continue
		}
		parsedPortHex, err := strconv.ParseInt(parts[1], 16, 32)
		if err != nil {
			continue
		}
		occupied[int(parsedPortHex)] = struct{}{}
	}

	// Faz o mesmo para IPv6 se existir
	if f6, err := os.Open("/proc/net/udp6"); err == nil {
		defer f6.Close()
		sc6 := bufio.NewScanner(f6)
		sc6.Scan()
		for sc6.Scan() {
			fields := strings.Fields(sc6.Text())
			if len(fields) < 2 {
				continue
			}
			parts := strings.SplitN(fields[1], ":", 2)
			if len(parts) != 2 {
				continue
			}
			parsedPortHex, err := strconv.ParseInt(parts[1], 16, 32)
			if err != nil {
				continue
			}
			occupied[int(parsedPortHex)] = struct{}{}
		}
	}

	for port := start; port <= end; port++ {
		if _, inUse := occupied[port]; !inUse {
			return port, true
		}
	}
	return 0, false
}

// findViaListen tenta abrir cada porta do range com ListenPacket (Windows :( ).
// Entradas: porta inicial e final do intervalo.
// Saída: primeira porta disponível ou erro quando nenhuma estiver livre.
func findViaListen(start, end int) (int, error) {
	for port := start; port <= end; port++ {
		conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue // ocupada ou sem permissão
		}
		conn.Close()
		return port, nil
	}
	return 0, logger.LogError(fmt.Sprintf(
		"NO_AVAILABLE_PORT: nenhuma porta disponível no intervalo %d-%d", start, end,
	))
}

// ExtractPort separa e retorna somente parte da porta de um endereço host:porta.
// Entradas: endereço de rede no formato host:porta.
// Saída: string da porta e erro se formato for inválido.
func ExtractPort(addr string) (string, error) {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", logger.LogError(fmt.Sprintf(
			"INVALID_ADDRESS: endereço inválido: %s", addr,
		))
	}
	return portStr, nil
}
