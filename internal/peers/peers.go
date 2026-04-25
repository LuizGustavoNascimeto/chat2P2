/*
Descrição: Implementa repositório concorrente de peers ativos e operações CRUD em memória.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package peers

import (
	"chat2p2/pkg/logger"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Peer struct {
	Addr     string
	LastSeen int64
	Waiting  bool
}

var (
	ErrInvalidPeer       = errors.New("invalid peer")
	ErrPeerAlreadyExists = errors.New("peer already exists")
	ErrPeerNotFound      = errors.New("peer not found")
)

type Store struct {
	mu    sync.RWMutex
	peers []Peer
}

// NewStore cria repositório vazio de peers.
// Entradas: nenhuma.
// Saída: ponteiro para Store inicializada.
func NewStore() *Store {
	return &Store{peers: make([]Peer, 0)}
}

// Create adiciona peer no repositório se endereço ainda não existir.
// Entradas: endereço do peer no formato host:porta.
// Saída: peer criado e erro quando endereço inválido ou duplicado.
func (s *Store) Create(addr string) (Peer, error) {
	if addr == "" {
		return Peer{}, ErrInvalidPeer
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.indexByAddr(addr) >= 0 {
		return Peer{}, ErrPeerAlreadyExists
	}

	log := logger.Get()
	log.Debug("Creating peer", zap.String("addr", addr))
	peer := Peer{
		Addr:     addr,
		Waiting:  false,
		LastSeen: time.Now().Unix(),
	}

	s.peers = append(s.peers, peer)
	return peer, nil
}

// Read busca peer pelo endereço.
// Entradas: endereço do peer.
// Saída: peer encontrado e erro quando não localizado.
func (s *Store) Read(addr string) (Peer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	index := s.indexByAddr(addr)
	if index < 0 {
		return Peer{}, ErrPeerNotFound
	}

	return s.peers[index], nil
}

// Update substitui peer existente por nova estrutura validada.
// Entradas: endereço original e objeto Peer atualizado.
// Saída: erro se peer não existir, endereço for inválido ou colidir com outro.
func (s *Store) Update(addr string, updatedPeer Peer) error {
	if addr == "" {
		return ErrInvalidPeer
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	index := s.indexByAddr(addr)
	if index < 0 {
		return ErrPeerNotFound
	}

	if updatedPeer.Addr == "" {
		updatedPeer.Addr = addr
	}

	if updatedPeer.Addr != addr && s.indexByAddr(updatedPeer.Addr) >= 0 {
		return ErrPeerAlreadyExists
	}

	s.peers[index] = updatedPeer
	return nil
}

// UpdateLastSeen atualiza timestamp de atividade do peer.
// Entradas: endereço do peer e timestamp Unix.
// Saída: erro quando peer não for encontrado.
func (s *Store) UpdateLastSeen(addr string, timestamp int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := s.indexByAddr(addr)
	if index < 0 {
		return ErrPeerNotFound
	}

	s.peers[index].LastSeen = timestamp
	return nil
}

// Delete remove peer pelo endereço.
// Entradas: endereço do peer.
// Saída: erro quando peer não existir.
func (s *Store) Delete(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := s.indexByAddr(addr)
	if index < 0 {
		return ErrPeerNotFound
	}

	s.peers = append(s.peers[:index], s.peers[index+1:]...)
	return nil
}

// List retorna cópia defensiva da lista de peers.
// Entradas: nenhuma.
// Saída: slice com todos peers armazenados.
func (s *Store) List() []Peer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Peer, len(s.peers))
	copy(result, s.peers)
	return result
}

// indexByAddr retorna índice interno de peer por endereço.
// Entradas: endereço do peer.
// Saída: índice >= 0 quando encontrado, ou -1 quando ausente.
func (s *Store) indexByAddr(addr string) int {
	for index, peer := range s.peers {
		if peer.Addr == addr {
			return index
		}
	}

	return -1
}

// GetAllAddresses retorna endereços de peers vistos recentemente.
// Entradas: nenhuma.
// Saída: slice de endereços filtrados por janela de 5 minutos.
func (s *Store) GetAllAddresses() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addresses := make([]string, len(s.peers))
	for index, peer := range s.peers {
		if peer.LastSeen < time.Now().Add(-5*time.Minute).Unix() {
			continue
		}
		addresses[index] = peer.Addr
	}

	return addresses
}
