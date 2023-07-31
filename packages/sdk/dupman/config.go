package dupman

import "time"

type Config struct {
	Debug       bool
	Timeout     time.Duration
	URL         string
	AccessToken string
}
