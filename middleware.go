package knock

import (
	"github.com/B9O2/knock/rawhttp"
	"github.com/B9O2/knock/rawhttp/client"
	"github.com/projectdiscovery/fastdialer/fastdialer"
)

type BaseMiddleware struct {
	f func(rawhttp.Options, fastdialer.Options, *client.Request)
}

func (bm *BaseMiddleware) Handle(opts rawhttp.Options, fdopts fastdialer.Options, req *client.Request) {
	bm.f(opts, fdopts, req)
}

func NewBaseMiddleware(f func(rawhttp.Options, fastdialer.Options, *client.Request)) *BaseMiddleware {
	return &BaseMiddleware{
		f: f,
	}
}
