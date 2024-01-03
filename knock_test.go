package knock

import (
	"fmt"
	"github.com/B9O2/Multitasking"
	"github.com/B9O2/knock/options"
	"github.com/B9O2/rawhttp"
	"github.com/B9O2/rawhttp/client"
	"github.com/projectdiscovery/fastdialer/fastdialer"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	k := NewClient()
	req := &BaseRequest{
		method:  POST,
		uri:     "/ok.php",
		headers: nil,
		version: HTTP_1_1,
		body:    nil,
	}
	for i := 0; i < 1; i++ {
		s, err := k.Knock("baidu.com", 443, false, req,
			//options.SetProxyOpt("http://127.0.0.1:8080", 5*time.Second),
			options.SetTimeoutOpt(15*time.Second),
			options.SetMiddlewareOpt("HelloWorld", NewBaseMiddleware(func(opts rawhttp.Options, fdopts fastdialer.Options, req *client.Request) {
				fmt.Println(req.Method, req.Headers, fdopts.Dialer.LocalAddr)
			})),
		)
		if err != nil {
			fmt.Println("fatal:", err)
		}
		fmt.Println(fmt.Sprintf("Connection: %s->%s by %s", s.LocalAddr(), s.RemoteAddr(), s.NetInterface().Name))
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
	req := &BaseRequest{
		method:  POST,
		uri:     "/word",
		headers: nil,
		version: HTTP_1_1,
		body:    nil,
	}

	mt := Multitasking.NewMultitasking("o", nil)
	mt.Register(func(dc Multitasking.DistributeController) {
		dc.Debug(true)
		for i := 0; i < 100; i++ {
			dc.AddTask(nil)
		}
	}, func(ec Multitasking.ExecuteController, a any) any {
		s, err := k.Knock("192.168.1.6", 8888, false, req,
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
		fmt.Println(fmt.Sprintf("Connection: %s->%s by %s", s.LocalAddr(), s.RemoteAddr(), s.NetInterface().Name))
		fmt.Println(resp.String())
		return i, nil
	}))
	mt.Run(50)

}
