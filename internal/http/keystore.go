package httpsrv

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// KeyStore holds the current install auth key and rotates it to a new random
// value once the configured TTL elapses. Old keys stop working immediately.
type KeyStore struct {
	mu        sync.Mutex
	key       string
	expiresAt time.Time
	ttl       time.Duration
}

func NewKeyStore(ttl time.Duration) *KeyStore {
	return &KeyStore{ttl: ttl}
}

// Current returns the active install key, generating a fresh one on first use
// and whenever the previous key has expired.
func (k *KeyStore) Current() string {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.key == "" || time.Now().After(k.expiresAt) {
		k.rotate()
	}
	return k.key
}

func (k *KeyStore) Valid(key string) bool {
	return key != "" && key == k.Current()
}

func (k *KeyStore) rotate() {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	k.key = hex.EncodeToString(b)
	k.expiresAt = time.Now().Add(k.ttl)
}
