package users

import (
	"chat2p2/pkg/logger"
	"errors"
	"sync"

	"go.uber.org/zap"
)

type User struct {
	Name     string
	Addr     string
	LastSeen int64
}

var (
	ErrInvalidUser       = errors.New("invalid user")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type Store struct {
	mu    sync.RWMutex
	users []User
}

func NewStore() *Store {
	return &Store{users: make([]User, 0)}
}

func (s *Store) Create(user User) error {
	if user.Addr == "" {
		return ErrInvalidUser
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.indexByAddr(user.Addr) >= 0 {
		return ErrUserAlreadyExists
	}
	log := logger.Get()
	log.Info("Creating user", zap.String("name", user.Name), zap.String("addr", user.Addr))

	s.users = append(s.users, user)
	return nil
}

func (s *Store) Read(addr string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx := s.indexByAddr(addr)
	if idx < 0 {
		return User{}, ErrUserNotFound
	}

	return s.users[idx], nil
}

func (s *Store) Update(addr string, updated User) error {
	if addr == "" {
		return ErrInvalidUser
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexByAddr(addr)
	if idx < 0 {
		return ErrUserNotFound
	}

	if updated.Addr == "" {
		updated.Addr = addr
	}

	if updated.Addr != addr && s.indexByAddr(updated.Addr) >= 0 {
		return ErrUserAlreadyExists
	}

	s.users[idx] = updated
	return nil
}

func (s *Store) UpdateLastSeen(addr string, timestamp int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := s.indexByAddr(addr)
	if idx < 0 {
		return ErrUserNotFound
	}
	s.users[idx].LastSeen = timestamp
	return nil
}

func (s *Store) Delete(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexByAddr(addr)
	if idx < 0 {
		return ErrUserNotFound
	}

	s.users = append(s.users[:idx], s.users[idx+1:]...)
	return nil
}

func (s *Store) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]User, len(s.users))
	copy(result, s.users)
	return result
}

func (s *Store) indexByAddr(addr string) int {
	for i, u := range s.users {
		if u.Addr == addr {
			return i
		}
	}

	return -1
}

func (s *Store) GetAllAddresses() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	addrs := make([]string, len(s.users))
	for i, u := range s.users {
		addrs[i] = u.Addr
	}
	return addrs

}
