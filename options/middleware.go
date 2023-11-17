package options

import (
	"github.com/B9O2/rawhttp"
)

type MiddlewareOpt struct {
	name string
	m    rawhttp.Middleware
}

func (mo MiddlewareOpt) Detail() (string, []string) {
	return "Middleware", []string{
		mo.name,
	}
}

func (mo MiddlewareOpt) Handle(opts *ClientOptions) error {
	opts.Middlewares = append(opts.Middlewares, mo.m)
	return nil
}

func SetMiddlewareOpt(name string, m rawhttp.Middleware) MiddlewareOpt {
	return MiddlewareOpt{
		name: name,
		m:    m,
	}
}
