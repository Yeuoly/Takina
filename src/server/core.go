package server

import (
	"encoding/json"
	"errors"

	"github.com/Yeuoly/Takina/src/proxy"
	"github.com/Yeuoly/Takina/src/types"
)

func server(request types.TakinaRequest) types.TakinaResponse {
	if request.Token != global_takina_instance.Token {
		return types.TakinaResponse{
			Code: -1,
			Msg:  "token error",
		}
	}

	switch request.Type {
	case types.TAKINA_TYPE_ADD_PROXY:
		var data types.TakinaRequestStartProxy
		err := json.Unmarshal([]byte(request.Data), &data)
		if err != nil {
			return types.TakinaResponse{
				Code: -1,
				Msg:  "failed to unmarshal",
			}
		}
		raddr, rport, err := addProxy(data.Laddr, data.Lport, data.ProxyType)
		if err != nil {
			return types.TakinaResponse{
				Code: -1,
				Msg:  err.Error(),
			}
		}
		response := types.TakinaResponseStartProxy{
			Raddr: raddr,
			Rport: rport,
		}
		text, _ := json.Marshal(response)
		return types.TakinaResponse{
			Code: 0,
			Msg:  "success",
			Data: string(text),
		}
	case types.TAKINA_TYPE_DEL_PROXY:
		var data types.TakinaRequestStopProxy
		err := json.Unmarshal([]byte(request.Data), &data)
		if err != nil {
			return types.TakinaResponse{
				Code: -1,
				Msg:  "failed to unmarshal",
			}
		}
		err = delProxy(data.Laddr, data.Lport)
		if err != nil {
			return types.TakinaResponse{
				Code: -1,
				Msg:  err.Error(),
			}
		}
		response := types.TakinaResponseStopProxy{}
		text, _ := json.Marshal(response)
		return types.TakinaResponse{
			Code: 0,
			Msg:  "success",
			Data: string(text),
		}
	}

	return types.TakinaResponse{
		Code: -1,
		Msg:  "unknown type",
	}
}

func addProxy(laddr string, lport int, protocol string) (string, int, error) {
	if protocol != "tcp" && protocol != "udp" && protocol != "http" && protocol != "https" && protocol != "stcp" && protocol != "sutp" {
		return "", 0, errors.New("unknown protocol")
	}
	return proxy.AutoLaunchProxy(laddr, lport, protocol)
}

func delProxy(laddr string, lport int) error {
	return proxy.StopProxy(laddr, lport)
}
