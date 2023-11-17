package options

type NetInterfaceOpt string

func (l NetInterfaceOpt) Detail() (string, []string) {
	return "NetInterface", []string{string(l)}
}

func (l NetInterfaceOpt) Handle(opts *ClientOptions) error {
	opts.NetInterface = string(l)
	return nil
}

func SetNetInterface(name string) NetInterfaceOpt {
	return NetInterfaceOpt(name)
}
