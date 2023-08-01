package dupman

import (
	"time"

	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
)

type Config struct {
	Credentials credentials.Provider
	Debug       bool
	Timeout     time.Duration
	URL         string
}
