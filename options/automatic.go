package options

import (
	"github.com/B9O2/knock/components"
)

type AutomaticOpt struct {
	HostHeader    bool
	ContentLength bool
}

func (a AutomaticOpt) Detail() (string, []string) {
	return "Automatic", []string{
		components.Ternary(a.HostHeader, "HostHeader: Enable", "HostHeader: Disable").(string),
		components.Ternary(a.ContentLength, "ContentLength: Enable", "ContentLength: Disable").(string),
	}
}

func (a AutomaticOpt) Handle(opts *ClientOptions) error {
	opts.AutomaticHostHeader = a.HostHeader
	opts.AutomaticContentLength = a.ContentLength
	return nil
}

func SetAutomaticOpt(hostHeader, contentLength bool) AutomaticOpt {
	return AutomaticOpt{
		HostHeader:    hostHeader,
		ContentLength: contentLength,
	}
}
