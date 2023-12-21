package options

import (
	"fmt"
	"net"
)

type LocalAddrOpt struct {
	addr *net.TCPAddr
	err  error
}

func (l LocalAddrOpt) Detail() (string, []string) {
	if l.err != nil {
		return "LocalAddr", []string{l.err.Error()}
	}
	return "LocalAddr", []string{l.addr.String()}
}

func (l LocalAddrOpt) Handle(opts *ClientOptions) error {
	opts.LocalAddr = l.addr
	return nil
}

func SetLocalAddr(host string, port int) LocalAddrOpt {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	return LocalAddrOpt{
		addr: addr,
		err:  err,
	}
}
