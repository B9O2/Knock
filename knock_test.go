package knock

import (
	"fmt"
	"github.com/B9O2/knock/options"
	"github.com/B9O2/rawhttp"
	"github.com/B9O2/rawhttp/client"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	k := NewClient()
	req := &BaseRequest{
		method:  POST,
		uri:     "/word",
		headers: nil,
		version: HTTP_1_0,
		body:    []byte("hello"),
	}
	s, err := k.Knock("n0p3.cn", 443, true, req,
		//options.SetProxyOpt("http://127.0.0.1:8080", 5*time.Second),
		options.SetTimeoutOpt(15*time.Second),
		options.SetMiddlewareOpt("HelloWorld", NewBaseMiddleware(func(opts rawhttp.Options, req *client.Request) {
			fmt.Println(req.Method, req.Headers, opts.FastDialerOpts.Dialer.LocalAddr)
		})),
	)
	if err != nil {
		fmt.Println("fatal:", err)
		return
	}
	fmt.Println(fmt.Sprintf("Connection: %s->%s by %s", s.LocalAddr(), s.RemoteAddr(), s.NetInterface().Name))
	resp, err := s.Response()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.String())
}
