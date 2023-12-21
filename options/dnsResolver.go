package options

type DNSListOpt []string

func (d DNSListOpt) Detail() (string, []string) {
	return "DNSList", d
}

func (d DNSListOpt) Handle(opts *ClientOptions) error {
	opts.BaseResolvers = d
	return nil
}

func SetDNSList(l []string) DNSListOpt {
	return l
}
