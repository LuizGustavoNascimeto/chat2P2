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

func NewStore() *Store {
	return &Store{peers: make([]Peer, 0)}
}

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
	log.Info("Creating peer", zap.String("addr", addr))
	peer := Peer{
		Addr:     addr,
		Waiting:  false,
		LastSeen: time.Now().Unix(),
	}

	s.peers = append(s.peers, peer)
	return peer, nil
}

func (s *Store) Read(addr string) (Peer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx := s.indexByAddr(addr)
	if idx < 0 {
		return Peer{}, ErrPeerNotFound
	}

	return s.peers[idx], nil
}

func (s *Store) Update(addr string, updated Peer) error {
	if addr == "" {
		return ErrInvalidPeer
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexByAddr(addr)
	if idx < 0 {
		return ErrPeerNotFound
	}

	if updated.Addr == "" {
		updated.Addr = addr
	}

	if updated.Addr != addr && s.indexByAddr(updated.Addr) >= 0 {
		return ErrPeerAlreadyExists
	}

	s.peers[idx] = updated
	return nil
}

func (s *Store) UpdateLastSeen(addr string, timestamp int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := s.indexByAddr(addr)
	if idx < 0 {
		return ErrPeerNotFound
	}
	s.peers[idx].LastSeen = timestamp
	return nil
}

func (s *Store) Delete(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexByAddr(addr)
	if idx < 0 {
		return ErrPeerNotFound
	}

	s.peers = append(s.peers[:idx], s.peers[idx+1:]...)
	return nil
}

func (s *Store) List() []Peer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Peer, len(s.peers))
	copy(result, s.peers)
	return result
}

func (s *Store) indexByAddr(addr string) int {
	for i, u := range s.peers {
		if u.Addr == addr {
			return i
		}
	}

	return -1
}

func (s *Store) GetAllAddresses() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	addrs := make([]string, len(s.peers))
	for i, u := range s.peers {
		addrs[i] = u.Addr
	}
	return addrs

}
