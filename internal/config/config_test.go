package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeConfig(t *testing.T, content string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_PATH", path)
}

func TestLoad(t *testing.T) {
	writeConfig(t, `
site:
  password: p
  install_key_ttl: 2h
server:
  addr: ":9090"
database:
  path: "/tmp/x.db"
`)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Site.Password != "p" {
		t.Errorf("password = %q, want %q", cfg.Site.Password, "p")
	}
	if cfg.Server.Addr != ":9090" {
		t.Errorf("addr = %q, want %q", cfg.Server.Addr, ":9090")
	}
	if cfg.Database.Path != "/tmp/x.db" {
		t.Errorf("db path = %q, want %q", cfg.Database.Path, "/tmp/x.db")
	}
	if got := cfg.InstallKeyTTL(); got != 2*time.Hour {
		t.Errorf("ttl = %v, want %v", got, 2*time.Hour)
	}
}

func TestLoadDefaults(t *testing.T) {
	writeConfig(t, `
site:
  password: p
`)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Addr != ":8080" {
		t.Errorf("default addr = %q, want :8080", cfg.Server.Addr)
	}
	if cfg.Database.Path != "app.db" {
		t.Errorf("default db path = %q, want app.db", cfg.Database.Path)
	}
	if got := cfg.InstallKeyTTL(); got != 24*time.Hour {
		t.Errorf("default ttl = %v, want %v", got, 24*time.Hour)
	}
}

func TestLoadMissingPassword(t *testing.T) {
	writeConfig(t, `
site:
  install_key_ttl: 1h
`)
	if _, err := Load(); err == nil {
		t.Fatal("expected error for missing site.password")
	}
}

func TestLoadBadTTL(t *testing.T) {
	writeConfig(t, `
site:
  password: p
  install_key_ttl: banana
`)
	if _, err := Load(); err == nil {
		t.Fatal("expected error for invalid install_key_ttl")
	}
}
