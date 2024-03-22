package knock

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"time"

	"github.com/B9O2/knock/components"
	"github.com/B9O2/knock/options"
	"github.com/B9O2/rawhttp"
	"github.com/B9O2/rawhttp/client"
	"github.com/projectdiscovery/fastdialer/fastdialer"
)

type Client struct {
	clientOpts rawhttp.Options
	opts       []options.Option
}

func (c *Client) parseOptions(opts ...options.Option) (*options.ClientOptions, error) {
	rawOpts := &options.ClientOptions{
		Options: c.clientOpts,
	}

	for _, opt := range opts {
		err := opt.Handle(rawOpts)
		if err != nil {
			return rawOpts, err
		}
	}

	return rawOpts, nil
}

func (c *Client) Knock(host string, port uint, https bool, req HTTPRequest, opts ...options.Option) (s *Snapshot, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

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
	sendOpts.Control = func(_, address string, c syscall.RawConn) (err error) {
		remoteAddr = address
		return nil
	}
	sendOpts.Middlewares = append(sendOpts.Middlewares, NewBaseMiddleware(func(opts rawhttp.Options, fdopts fastdialer.Options, req *client.Request) {
		if localAddr := fdopts.Dialer.LocalAddr; localAddr != nil {
			s.ci.localAddr = append(s.ci.localAddr, localAddr.(*net.TCPAddr))
		}
	}))

	ct := rawhttp.NewClient(sendOpts.Options)
	defer ct.Close()
	//send
	var reader *bytes.Buffer
	if req.Body() != nil {
		reader = bytes.NewBuffer(req.Body())
	} else {
		reader = bytes.NewBuffer([]byte{})
	}
	resp, connErr := ct.DoRawWithOptions(
		string(req.Method()),
		targetURL,
		req.URI(),
		client.Version(req.Version()),
		req.Headers(),
		reader,
		sendOpts.Options,
	)

	// TODO 移除
	fmt.Println(targetURL, sendOpts)
	fmt.Println(string(req.Raw()))

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
		return s, connErr
	}

	//Response
	body, err := ReadAll(resp.Body)
	if err != nil {
		s.ci.log("<Knock::ReadBody> ", err.Error())
		s.ci.err = errors.New("<Knock::ReadBody> " + err.Error())
	}
	s.resp = &Response{
		resp,
		body,
	}

	return s, nil
}

func ReadAll(r io.Reader) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				fmt.Println("!!!!!!", string(b), err)
			}
			return b, err
		}

		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
	}
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
	}
	c := Client{
		clientOpts: rawHTTPOpts,
		opts:       opts,
	}
	return &c
}
