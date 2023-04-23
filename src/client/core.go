package server

import (
	"errors"

	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/proxy"
	"github.com/Yeuoly/Takina/src/types"
)

func addProxy(laddr string, lport int, protocol string) (string, int, error) {
	if protocol != "tcp" && protocol != "udp" && protocol != "http" && protocol != "https" && protocol != "stcp" && protocol != "sutp" {
		return "", 0, errors.New("unknown protocol")
	}

	// request a port in random node
	node, err := GetTakina().RandomNode()
	if err != nil {
		return "", 0, err
	}

	takina := GetTakina()
	resp, err := helper.SendPostAndParse[types.TakinaResponseWarp[types.TakinaResponseGetPort]](
		node.GenerateUrl(types.ROUTER_TAKINA_SERVER_REQUEST_PORT),
		helper.HttpPayloadJson(GetPackedRequest(takina, types.TakinaRequestGetPort{})),
		helper.HttpTimeout(2000),
	)

	if err != nil {
		return "", 0, err
	}

	if resp.Error() != nil {
		return "", 0, resp.Error()
	}

	return proxy.AutoLaunchProxy(laddr, lport, protocol, node.Address, resp.Data.Port)
}

func delProxy(laddr string, lport int) error {
	return proxy.StopProxy(laddr, lport)
}

func listProxy() []proxy.Proxy {
	return proxy.GetProxies()
}
