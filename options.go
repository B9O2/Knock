package knock

import "time"

type KnockOptions struct {
	Timeout       time.Duration
	HTTPProxyAddr string
}
