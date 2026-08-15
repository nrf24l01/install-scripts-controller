package httpsrv

import (
	"testing"
	"time"
)

func TestKeyStoreRotation(t *testing.T) {
	ks := NewKeyStore(50 * time.Millisecond)

	k1 := ks.Current()
	if k1 == "" {
		t.Fatal("expected a generated key")
	}
	if !ks.Valid(k1) {
		t.Fatal("fresh key should be valid")
	}

	// Same key until it expires.
	if k2 := ks.Current(); k2 != k1 {
		t.Errorf("key changed before expiry: %q != %q", k1, k2)
	}

	time.Sleep(120 * time.Millisecond)

	if ks.Valid(k1) {
		t.Error("expired key should no longer be valid")
	}

	k3 := ks.Current()
	if k3 == k1 {
		t.Error("expected a new key after expiry")
	}
	if !ks.Valid(k3) {
		t.Error("new key should be valid")
	}
}

func TestKeyStoreEmptyInvalid(t *testing.T) {
	ks := NewKeyStore(time.Hour)
	if ks.Valid("") {
		t.Error("empty key should be invalid")
	}
}
