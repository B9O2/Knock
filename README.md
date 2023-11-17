# Knock

- 即使TCP连接失败也可以获取LocalAddr与RemoteAddr
- 重定向可以拿到过程中的所有LocalAddr
- 可以指定网络接口(Net Interface)
- 支持中间件，可以在写入TCP连接前进行操作

```go
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
```

