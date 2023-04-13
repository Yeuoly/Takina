package server

import (
	"errors"

	"github.com/Yeuoly/Takina/src/proxy"
)

func addProxy(laddr string, lport int, protocol string) (string, int, error) {
	if protocol != "tcp" && protocol != "udp" && protocol != "http" && protocol != "https" && protocol != "stcp" && protocol != "sutp" {
		return "", 0, errors.New("unknown protocol")
	}
	return proxy.AutoLaunchProxy(laddr, lport, protocol)
}

func delProxy(laddr string, lport int) error {
	return proxy.StopProxy(laddr, lport)
}

func listProxy() []proxy.Proxy {
	return proxy.GetProxies()
}