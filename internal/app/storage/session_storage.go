package storage

import (
	"sync"
	"time"
)

// SessionStorage данные по авторизованным пользователям
type SessionStorage struct {
	values map[string]time.Time
	mx     sync.RWMutex
}

type SessionManager interface {
	Add(token string, expire time.Time)
	IsValid(token string) bool
}

func NewSessionStorage() SessionManager {
	return &SessionStorage{
		values: make(map[string]time.Time),
	}
}

func (s *SessionStorage) Add(token string, expire time.Time) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.values[token] = expire
}

func (s *SessionStorage) IsValid(token string) bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	expireTime, ok := s.values[token]
	if !ok {
		return false
	}

	return expireTime.After(time.Now())
}
