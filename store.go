package sessions_you_dream_about

import (
	"errors"
	"sync"
	"time"
)

type Store interface {
	Save(token string, session SessionData) error
	Get(token string) (SessionData, error)
	Delete(token string) error
	Cleanup(ttl time.Duration) error
}

type InMemoryStore struct {
	sessions map[string]SessionData
	mutex    sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		sessions: make(map[string]SessionData),
	}
}

func (s *InMemoryStore) Save(token string, session SessionData) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions[token] = session
	return nil
}

func (s *InMemoryStore) Get(token string) (SessionData, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	session, exists := s.sessions[token]
	if !exists {
		return SessionData{}, errors.New("session not found")
	}
	return session, nil
}

func (s *InMemoryStore) Delete(token string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, token)
	return nil
}

func (s *InMemoryStore) Cleanup(ttl time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	now := time.Now()
	for token, session := range s.sessions {
		if now.Sub(session.LastActivity) > ttl {
			delete(s.sessions, token)
		}
	}
	return nil
}