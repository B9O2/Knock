package options

import (
	"fmt"

	"github.com/B9O2/rawhttp"
)

type Option interface {
	Detail() (string, []string)
	Handle(opts *ClientOptions) error
}

type ClientOptions struct {
	rawhttp.Options
}

func (co ClientOptions) String() string {
	return fmt.Sprintf(`Timeout %s 
	FollowRedirects        %t
    MaxRedirects           %d
    AutomaticHostHeader    %t
    AutomaticContentLength %t
    CustomHeaders          %v
    ForceReadAllBody       %t
    CustomRawBytes         %s
    Proxy                  %s
    ProxyDialTimeout       %s
    SNI                    %s`,
		co.Timeout.String(),
		co.FollowRedirects,
		co.MaxRedirects,
		co.AutomaticHostHeader,
		co.AutomaticContentLength,
		co.CustomHeaders,
		co.ForceReadAllBody,
		co.CustomRawBytes,
		co.Proxy,
		co.ProxyDialTimeout,
		co.SNI,
	)
}
