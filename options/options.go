package options

import "github.com/B9O2/rawhttp"

type Option interface {
	Detail() (string, []string)
	Handle(opts *ClientOptions) error
}

type ClientOptions struct {
	rawhttp.Options
}
