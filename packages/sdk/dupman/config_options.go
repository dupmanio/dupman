package dupman

import (
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
)

type Option func(config *Config)

// WithCredentials defines Credentials Provider for dupman Config.
func WithCredentials(credentials credentials.Provider) Option {
	return func(config *Config) {
		config.Credentials = credentials
	}
}

// WithBaseURL defines Base URL for dupman Config.
func WithBaseURL(url string) Option {
	return func(config *Config) {
		config.BaseURL = url
	}
}

// WithDebug defines Debug option for dupman Config.
func WithDebug(debug bool) Option {
	return func(config *Config) {
		config.Debug = debug
	}
}
