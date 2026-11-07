package config

import (
	"os"
	"strconv"
)

// Config holds runtime configuration options for ogcraft.
type Config struct {
	Port         int
	Host         string
	CacheDir     string
	TemplatesDir string
	L1MaxBytes   int64
	Secret       string
	RateLimit    int
}

// DefaultConfig returns reasonable defaults suitable for zero-config startup.
func DefaultConfig() *Config {
	return &Config{
		Port:         8080,
		Host:         "0.0.0.0",
		CacheDir:     "./data/cache",
		TemplatesDir: "./data/templates",
		L1MaxBytes:   24 * 1024 * 1024, // 24MB
		Secret:       "",
		RateLimit:    100,
	}
}

// LoadFromEnv loads environment variable overrides onto the configuration.
func (c *Config) LoadFromEnv() {
	if val := os.Getenv("OGCRAFT_PORT"); val != "" {
		if p, err := strconv.Atoi(val); err == nil && p > 0 {
			c.Port = p
		}
	}
	if val := os.Getenv("OGCRAFT_HOST"); val != "" {
		c.Host = val
	}
	if val := os.Getenv("OGCRAFT_CACHE_DIR"); val != "" {
		c.CacheDir = val
	}
	if val := os.Getenv("OGCRAFT_TEMPLATES_DIR"); val != "" {
		c.TemplatesDir = val
	}
	if val := os.Getenv("OGCRAFT_L1_BYTES"); val != "" {
		if b, err := strconv.ParseInt(val, 10, 64); err == nil && b > 0 {
			c.L1MaxBytes = b
		}
	}
	if val := os.Getenv("OGCRAFT_SECRET"); val != "" {
		c.Secret = val
	}
	if val := os.Getenv("OGCRAFT_RATE_LIMIT"); val != "" {
		if r, err := strconv.Atoi(val); err == nil && r > 0 {
			c.RateLimit = r
		}
	}
}
