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

// FindAvailablePortUDP retorna a primeira porta UDP livre no intervalo [start, end].
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

// findViaProc lê /proc/net/udp (Linux) e retorna a primeira porta do range
// que não constar como ocupada. Muito mais rápido que tentar bind em cada porta.
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
		port64, err := strconv.ParseInt(parts[1], 16, 32)
		if err != nil {
			continue
		}
		occupied[int(port64)] = struct{}{}
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
			port64, err := strconv.ParseInt(parts[1], 16, 32)
			if err != nil {
				continue
			}
			occupied[int(port64)] = struct{}{}
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

func ExtractPort(addr string) (string, error) {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", logger.LogError(fmt.Sprintf(
			"INVALID_ADDRESS: endereço inválido: %s", addr,
		))
	}
	return portStr, nil
}
