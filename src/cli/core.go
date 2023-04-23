package cli

import (
	"fmt"

	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
)

func startProxy(laddr string, lport int, typ string) (types.Proxy, error) {
	resp, err := helper.SendPostAndParse[types.TakinaResponseWarp[types.TakinaResponseStartProxy]](
		fmt.Sprintf("http://127.0.0.1:%d%s", *takina_port, types.ROUTER_TAKINA_CLIENT_DAEMON_START_PROXY),
		helper.HttpPayloadJson(
			generateRequest(types.TakinaRequestStartProxy{
				Laddr:     laddr,
				Lport:     lport,
				ProxyType: typ,
			}),
		),
	)
	if err != nil {
		return types.Proxy{}, err
	}
	if resp.Error() != nil {
		return types.Proxy{}, err
	}
	return types.Proxy{
		Laddr: laddr,
		Lport: lport,
		Raddr: resp.Data.Raddr,
		Rport: resp.Data.Rport,
		Type:  typ,
	}, nil
}

func stopProxy(laddr string, lport int) error {
	resp, err := helper.SendPostAndParse[types.TakinaResponseWarp[types.TakinaResponseStopProxy]](
		fmt.Sprintf("http://127.0.0.1:%d%s", *takina_port, types.ROUTER_TAKINA_CLIENT_DAEMON_STOP_PROXY),
		helper.HttpPayloadJson(
			generateRequest(types.TakinaRequestStopProxy{
				Laddr: laddr,
				Lport: lport,
			}),
		),
	)
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return err
	}
	return nil
}

func listProxy() ([]types.Proxy, error) {
	resp, err := helper.SendGetAndParse[types.TakinaResponseWarp[types.TakinaResponseListProxy]](
		fmt.Sprintf("http://127.0.0.1:%d%s", *takina_port, types.ROUTER_TAKINA_CLIENT_DAEMON_LIST_PROXY),
		helper.HttpPayloadJson(generateRequest(types.TakinaRequestListProxy{})),
	)
	if err != nil {
		return nil, err
	}

	if resp.Error() != nil {
		return nil, err
	}

	return resp.Data.Proxies, nil
}
