package knock

import (
	"github.com/B9O2/rawhttp"
	"github.com/B9O2/rawhttp/client"
)

type BaseMiddleware struct {
	f func(rawhttp.Options, *client.Request)
}

func (bm *BaseMiddleware) Handle(opts rawhttp.Options, req *client.Request) {
	bm.f(opts, req)
}

func NewBaseMiddleware(f func(rawhttp.Options, *client.Request)) *BaseMiddleware {
	return &BaseMiddleware{
		f: f,
	}
}
