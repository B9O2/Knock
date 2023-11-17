package knock

import (
	"github.com/B9O2/rawhttp"
	"github.com/B9O2/rawhttp/client"
)

type BaseMiddleware struct {
	f func(rawhttp.Options, *client.Request, rawhttp.Conn)
}

func (bm *BaseMiddleware) Handle(opts rawhttp.Options, req *client.Request, conn rawhttp.Conn) {
	bm.f(opts, req, conn)
}

func NewBaseMiddleware(f func(rawhttp.Options, *client.Request, rawhttp.Conn)) *BaseMiddleware {
	return &BaseMiddleware{
		f: f,
	}
}
