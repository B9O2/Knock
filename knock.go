package knock

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"time"

	"github.com/B9O2/knock/components"
	"github.com/B9O2/knock/rawhttp"
	"github.com/projectdiscovery/fastdialer/fastdialer"
)

type Client struct {
	defaultOpts *rawhttp.Options
}

func (c *Client) SetDefaultOptions(opts *rawhttp.Options) {
	if opts != nil {
		c.defaultOpts = opts
	} else {
		c.defaultOpts = rawhttp.DefaultOptions
	}
}

func (c *Client) DefaultOptions() rawhttp.Options {
	return *c.defaultOpts
}

func (c *Client) Knock(host string, port uint, https bool, req HTTPRequest, opts *KnockOptions) (s *Snapshot, err error) {
	var deadline time.Time
	if opts.Timeout != 0 {
		deadline = time.Now().Add(opts.Timeout)
	}

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
	clentOpts := c.defaultOpts

	//dialer setting
	remoteAddr := ""
	dialerOpts := fastdialer.DefaultOptions

	dialerOpts.Dialer = &net.Dialer{
		Deadline: deadline,
		ControlContext: func(ctx context.Context, network, address string, c syscall.RawConn) error {
			remoteAddr = address
			return nil
		},
	}
	dialer, err := fastdialer.NewDialer(dialerOpts)
	clentOpts.FastDialer = dialer

	if len(opts.HTTPProxyAddr) > 0 {
		clentOpts.Proxy = opts.HTTPProxyAddr
	}
	if int64(opts.Timeout) > 0 {
		clentOpts.ProxyDialTimeout = opts.Timeout
	}

	ct := rawhttp.NewClient(clentOpts)
	defer ct.Close()
	//send
	var reader *bytes.Buffer
	if req.Body() != nil {
		reader = bytes.NewBuffer(req.Body())
	} else {
		reader = bytes.NewBuffer([]byte{})
	}
	resp, connErr := ct.DoRaw(
		string(req.Method()),
		targetURL,
		req.URI(),
		req.Headers(),
		reader,
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
		return s, connErr
	}

	//Response
	body, err := io.ReadAll(resp.Body)
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

func NewClient() *Client {
	// rawHTTPOpts := rawhttp.Options{
	// 	Timeout:                5 * time.Second,
	// 	FollowRedirects:        true,
	// 	MaxRedirects:           10,
	// 	AutomaticHostHeader:    true,
	// 	AutomaticContentLength: true,
	// 	CustomHeaders:          nil,
	// 	ForceReadAllBody:       false,
	// 	CustomRawBytes:         nil,
	// 	Proxy:                  "",
	// 	ProxyDialTimeout:       5 * time.Second,
	// 	SNI:                    "",
	// }
	c := Client{
		defaultOpts: rawhttp.DefaultOptions,
	}
	return &c
}
