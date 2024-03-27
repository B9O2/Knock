package knock

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/B9O2/Multitasking"
	"github.com/B9O2/knock/options"
	"github.com/B9O2/rawhttp"
	"github.com/B9O2/rawhttp/client"
	"github.com/projectdiscovery/fastdialer/fastdialer"
)

func TestNewClient(t *testing.T) {
	k := NewClient(options.SetProxyOpt("http://127.0.0.1:8083", 5*time.Second))
	req := HTTPRequest{
		method:  GET,
		uri:     "/users/sign_in",
		headers: map[string][]string{},
		version: HTTP_1_1,
	}
	for i := 0; i < 1; i++ {
		s, err := k.Knock("192.168.31.98", 8080, false, req,
			options.SetProxyOpt("http://127.0.0.1:8083", 5*time.Second),
			options.SetTimeoutOpt(15*time.Second),
			options.SetMiddlewareOpt("HelloWorld", NewBaseMiddleware(func(opts rawhttp.Options, fdopts fastdialer.Options, req *client.Request) {
				fmt.Println(req.Method, req.Headers, fdopts.Dialer.LocalAddr)
			})),
		)
		if err != nil {
			fmt.Println("fatal:", err)
		}
		fmt.Printf("Connection: %s->%s by %s\n", s.LocalAddr(), s.RemoteAddr(), s.NetInterface().Name)
		resp, err := s.Response()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp.String())
	}

}

func TestNewMultitaskingClient(t *testing.T) {
	file, err := os.Create("./cpu.pprof")
	if err != nil {
		fmt.Printf("create cpu pprof failed, err:%v\n", err)
		return
	}
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	k := NewClient()
	req := HTTPRequest{
		method: POST,
		uri:    "/index.php?s=captcha",
		headers: map[string][]string{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
		version: HTTP_1_1,
		body:    []byte("_method=__construct&filter[]=printf&method=get&server[REQUEST_METHOD]=uhoaycxtvnkhybm"),
	}

	mt := Multitasking.NewMultitasking("o", nil)
	mt.Register(func(dc Multitasking.DistributeController) {
		dc.Debug(true)
		for i := 0; i < 100; i++ {
			dc.AddTask(nil)
		}
	}, func(ec Multitasking.ExecuteController, a any) any {
		s, err := k.Knock("192.168.1.13", 8080, false, req,
			//options.SetProxyOpt("http://127.0.0.1:8081", 5*time.Second),
			options.SetTimeoutOpt(5*time.Second),
		)
		if err != nil {
			fmt.Println("fatal:", err)
			return err
		}

		return s
	})
	mt.SetResultMiddlewares(Multitasking.NewBaseMiddleware(func(ec Multitasking.ExecuteController, i interface{}) (interface{}, error) {
		s, ok := i.(*Snapshot)
		if !ok {
			fmt.Println(i.(error))
			return s, i.(error)
		}
		resp, err := s.Response()
		if err != nil {
			fmt.Println(err)
			return s, err
		}
		fmt.Printf("Connection: %s->%s by %s\n", s.LocalAddr(), s.RemoteAddr(), s.NetInterface().Name)
		fmt.Println(resp.String())
		return i, nil
	}))
	mt.Run(context.Background(), 50)

}
