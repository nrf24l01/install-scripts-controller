package httpsrv

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type SessionStore struct {
	mu     sync.Mutex
	tokens map[string]time.Time
	ttl    time.Duration
}

func NewSessionStore(ttl time.Duration) *SessionStore {
	s := &SessionStore{
		tokens: make(map[string]time.Time),
		ttl:    ttl,
	}
	go s.cleanupLoop()
	return s
}

func (s *SessionStore) Create() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	tok := hex.EncodeToString(b)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[tok] = time.Now().Add(s.ttl)
	return tok
}

func (s *SessionStore) Valid(tok string) bool {
	if tok == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.tokens[tok]
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		delete(s.tokens, tok)
		return false
	}
	return true
}

func (s *SessionStore) Delete(tok string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, tok)
}

func (s *SessionStore) cleanupLoop() {
	t := time.NewTicker(s.ttl / 2)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		s.mu.Lock()
		for k, exp := range s.tokens {
			if now.After(exp) {
				delete(s.tokens, k)
			}
		}
		s.mu.Unlock()
	}
}
