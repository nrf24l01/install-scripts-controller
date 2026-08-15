package httpsrv

import (
	"testing"
	"time"
)

func TestSessionStoreLifecycle(t *testing.T) {
	ss := NewSessionStore(50 * time.Millisecond)

	tok := ss.Create()
	if tok == "" {
		t.Fatal("expected a token")
	}
	if !ss.Valid(tok) {
		t.Fatal("fresh token should be valid")
	}

	ss.Delete(tok)
	if ss.Valid(tok) {
		t.Fatal("deleted token should be invalid")
	}

	tok2 := ss.Create()
	time.Sleep(120 * time.Millisecond)
	if ss.Valid(tok2) {
		t.Fatal("expired token should be invalid")
	}
}
