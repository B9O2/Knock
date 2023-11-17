package components

import (
	"errors"
	"net"
)

func QueryNetInterface(ip net.IP) (net.Interface, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return net.Interface{}, err
	}

	for _, iface := range netInterfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ipNet.Contains(ip) {
				return iface, nil
			}
		}
	}

	return net.Interface{
		Name: "Unknown(interface)",
	}, errors.New("no net interface matched with '" + ip.String() + "'")
}
