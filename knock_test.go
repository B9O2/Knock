package knock

import (
	"Knock/options"
	"fmt"
	"github.com/B9O2/rawhttp"
	"github.com/B9O2/rawhttp/client"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	/*
		u, err := url.Parse("socks5://127.0.0.1:8080")
		if err != nil {
			fmt.Println("!!!")
		} else {
			fmt.Println(u.Scheme)
		}
		return
	*/
	k := NewClient()
	req := &BaseRequest{
		method:  GET,
		uri:     "",
		headers: nil,
		body:    nil,
	}
	s, err := k.Knock("192.168.1.14", 81, false, req,
		options.SetProxyOpt("http://127.0.0.1:8080", 1*time.Second),
		options.SetMiddlewareOpt("HelloWorld", NewBaseMiddleware(func(opts rawhttp.Options, req *client.Request, conn rawhttp.Conn) {
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
