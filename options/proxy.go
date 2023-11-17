package options

import (
	"time"
)

type ProxyOpt struct {
	addr    string
	timeout time.Duration
}

func (p ProxyOpt) Detail() (string, []string) {
	return "Proxy", []string{
		"Addr: " + p.addr,
		"Timeout: " + p.timeout.String(),
	}
}

func (p ProxyOpt) Handle(opts *ClientOptions) error {
	opts.Proxy = p.addr
	opts.ProxyDialTimeout = p.timeout
	return nil
}

func SetProxyOpt(addr string, timeout time.Duration) ProxyOpt {
	return ProxyOpt{
		addr:    addr,
		timeout: timeout,
	}
}
