package dupman

import (
	"time"

	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
)

const (
	defaultTimeout    = 30 * time.Second
	defaultRetryCount = 3
)

type Config struct {
	Credentials credentials.Provider
	Debug       bool
	Timeout     time.Duration
	RetryCount  int
	BaseURL     string
}

// NewConfig creates a dupman reusable dupman service configuration.
func NewConfig(options ...Option) *Config {
	config := &Config{
		Timeout:    defaultTimeout,
		RetryCount: defaultRetryCount,
	}

	for _, opt := range options {
		opt(config)
	}

	return config
}
