package knock

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/B9O2/Multitasking"
)

func TestNewClient(t *testing.T) {
	k := NewClient()
	req := HTTPRequest{
		method:  GET,
		uri:     "/users/sign_in",
		headers: map[string][]string{},
		version: HTTP_1_1,
	}
	for i := 0; i < 1; i++ {
		s, err := k.Knock("127.0.0.1", 8080, false, req, &KnockOptions{
			Timeout: time.Second * time.Duration(3),
		})
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

	mt := Multitasking.NewMultitasking("o", nil)
	mt.Register(func(dc Multitasking.DistributeController) {
		dc.Debug(true)
		for i := 0; i < 20; i++ {
			dc.AddTask(nil)
		}
	}, func(ec Multitasking.ExecuteController, a any) any {
		req := HTTPRequest{
			method: POST,
			uri:    "/index.php?s=captcha",
			headers: map[string][]string{
				"Content-Type": {"application/x-www-form-urlencoded"},
			},
			version: HTTP_1_1,
			body:    []byte("_method=__construct&filter[]=printf&method=get&server[REQUEST_METHOD]=uhoaycxtvnkhybm"),
		}
		s, err := k.Knock("baidu.com", 443, true, req, &KnockOptions{
			Timeout: time.Second * time.Duration(3),
		})
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
	mt.Run(context.Background(), 20)

}
