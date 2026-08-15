package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Site     SiteConfig     `yaml:"site"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`

	installKeyTTL time.Duration
}

type SiteConfig struct {
	Password      string `yaml:"password"`
	InstallKeyTTL string `yaml:"install_key_ttl"`
	PublicURL     string `yaml:"public_url"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// InstallKeyTTL returns how long a randomly generated install key stays valid.
func (c *Config) InstallKeyTTL() time.Duration {
	return c.installKeyTTL
}

func Default() *Config {
	return &Config{
		Server:   ServerConfig{Addr: ":8080"},
		Database: DatabaseConfig{Path: "app.db"},
	}
}

// Load reads the config from the path given by the CONFIG_PATH env var,
// falling back to config.yml in the working directory.
func Load() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.yml"
	}

	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":8080"
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path = "app.db"
	}
	if cfg.Site.Password == "" {
		return nil, fmt.Errorf("config %s: site.password is required", path)
	}
	if cfg.Site.InstallKeyTTL == "" {
		cfg.Site.InstallKeyTTL = "24h"
	}
	ttl, err := time.ParseDuration(cfg.Site.InstallKeyTTL)
	if err != nil {
		return nil, fmt.Errorf("config %s: site.install_key_ttl is invalid: %w", path, err)
	}
	cfg.installKeyTTL = ttl

	return cfg, nil
}
