package options

import (
	"fmt"
)

type DNSListOpt []string

func (d DNSListOpt) Detail() (string, []string) {
	return "DNSList", d
}

func (d DNSListOpt) Handle(opts *ClientOptions) error {
	fmt.Println(opts.FastDialerOpts.BaseResolvers)
	opts.FastDialerOpts.BaseResolvers = d
	return nil
}

func SetDNSList(l []string) DNSListOpt {
	return l
}
