package main

import (
	"fmt"
	"time"

	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
)

func generateRequest[T any](data T) types.TakinaRequest[T] {
	return types.TakinaRequest[T]{
		Token: "testtest",
		Data:  data,
	}
}

func startProxy(laddr string, lport int, typ string) (types.Proxy, error) {
	resp, err := helper.SendPostAndParse[types.TakinaResponseWarp[types.TakinaResponseStartProxy]](
		fmt.Sprintf("http://127.0.0.1:%d%s", 40002, types.ROUTER_TAKINA_CLIENT_DAEMON_START_PROXY),
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

func main() {
	for i := 80; i < 20080; i++ {
		proxy, err := startProxy("127.0.0.1", 8081, "tcp")
		if err != nil {
			fmt.Println(err)
		}
		j := 0
		for ; j < 3; j++ {
			resp, err := helper.HttpGet(fmt.Sprintf("http://127.0.0.1:%d", proxy.Rport), make(map[string]string))
			if len(resp) < 20 || err != nil {
				fmt.Println(resp, err)
			} else {
				fmt.Println("ok")
				break
			}
			time.Sleep(time.Millisecond * 50)
		}
		if j == 3 {
			panic("failed")
		}
	}
}