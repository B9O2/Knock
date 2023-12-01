package knock

import (
	"bytes"
	"fmt"
	"github.com/B9O2/knock/components"
	"github.com/B9O2/knock/options"
	"github.com/B9O2/rawhttp"
	"github.com/B9O2/rawhttp/client"
	"github.com/projectdiscovery/fastdialer/fastdialer"
	"net"
	"syscall"
	"time"
)

type Client struct {
	clientOpts rawhttp.Options
	opts       []options.Option
}

func (c *Client) parseOptions(opts ...options.Option) (*options.ClientOptions, error) {
	rawOpts := &options.ClientOptions{
		Options: &c.clientOpts,
	}
	rawOpts.FastDialerOpts.Dialer = &net.Dialer{
		Timeout:   rawOpts.FastDialerOpts.DialerTimeout,
		KeepAlive: rawOpts.FastDialerOpts.DialerKeepAlive,
		DualStack: true,
	}

	for _, opt := range opts {
		err := opt.Handle(rawOpts)
		if err != nil {
			return rawOpts, err
		}
	}

	return rawOpts, nil
}

func (c *Client) Knock(host string, port uint, https bool, req Request, opts ...options.Option) (s *Snapshot, err error) {
	s = &Snapshot{
		req: req,
		ci: &ConnectionInfo{
			events: make([]Event, 0),
		},
	}

	protocol := "http"
	if https {
		protocol = "https"
	}

	targetURL := fmt.Sprintf("%s://%s:%d", protocol, host, port)
	sendOpts, err := c.parseOptions(append(c.opts, opts...)...)
	if err != nil {
		return
	}

	//dialer setting
	remoteAddr := ""
	sendOpts.FastDialerOpts.Dialer.Control = func(_, address string, c syscall.RawConn) (err error) {
		remoteAddr = address
		return nil
	}
	sendOpts.Middlewares = append(sendOpts.Middlewares, NewBaseMiddleware(func(opts rawhttp.Options, req *client.Request) {
		s.ci.localAddr = append(s.ci.localAddr, opts.FastDialerOpts.Dialer.LocalAddr.(*net.TCPAddr))
	}))

	ct := rawhttp.NewClient(sendOpts.Options)
	defer ct.Close()

	//send
	resp, connErr := ct.DoRawWithOptions(
		string(req.Method()),
		targetURL,
		req.URI(),
		req.Headers(),
		bytes.NewReader(req.Body()),
		sendOpts.Options,
	)

	//after request
	var terr error
	if s.ci.remoteAddr, terr = net.ResolveTCPAddr("tcp", remoteAddr); terr != nil {
		s.ci.log("ConnectionInfo::RemoteAddr", terr.Error())
	}
	if len(s.ci.localAddr) > 0 {
		if s.ci.inter, terr = components.QueryNetInterface(s.ci.localAddr[0].IP); terr != nil {
			s.ci.log("ConnectionInfo::NetInterface", terr.Error())
		}
	}

	if connErr != nil {
		s.ci.log("Knock", connErr.Error())
		s.ci.err = connErr
		return
	}

	//Response
	s.resp = &Response{
		resp,
	}

	return
}

func NewClient(opts ...options.Option) *Client {
	rawHTTPOpts := rawhttp.Options{
		Timeout:                5 * time.Second,
		FollowRedirects:        true,
		MaxRedirects:           10,
		AutomaticHostHeader:    true,
		AutomaticContentLength: true,
		CustomHeaders:          nil,
		ForceReadAllBody:       false,
		CustomRawBytes:         nil,
		Proxy:                  "",
		ProxyDialTimeout:       5 * time.Second,
		SNI:                    "",
		FastDialerOpts:         &fastdialer.DefaultOptions,
	}
	c := Client{
		clientOpts: rawHTTPOpts,
		opts:       opts,
	}
	return &c
}
