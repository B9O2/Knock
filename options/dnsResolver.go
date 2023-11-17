package options

import (
	"fmt"
	"github.com/B9O2/rawhttp"
	"github.com/projectdiscovery/fastdialer/fastdialer"
)

type DNSListOpt []string

func (d DNSListOpt) Detail() (string, []string) {
	return "DNSList", d
}

func (d DNSListOpt) Handle(opts *rawhttp.Options, fopts *fastdialer.Options) error {
	fmt.Println(fopts.BaseResolvers)
	fopts.BaseResolvers = d
	return nil
}

func SetDNSList(l []string) DNSListOpt {
	return l
}
