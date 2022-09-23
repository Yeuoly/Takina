package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Yeuoly/Takina/src/proxy"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"gopkg.in/yaml.v2"
)

type Takina struct {
	Token string `yaml:"token"`
}

const MESSAGE_TAKINA = 0x0001

type TakinaServer struct {
	znet.BaseRouter
}

var global Takina

func init() {
	global = Takina{
		Token: "takina",
	}

	data, err := ioutil.ReadFile("conf/server.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &global)
	if err != nil {
		panic(err)
	}
}

func (root Takina) Run() {
	server := znet.NewServer()
	server.AddRouter(MESSAGE_TAKINA, &TakinaServer{})
	server.Serve()
}

func GetTakina() Takina {
	return global
}

func (router *TakinaServer) Handle(req ziface.IRequest) {
	var request TakinaRequest
	//unmarshal
	err := json.Unmarshal(req.GetData(), &request)
	if err != nil {
		req.GetConnection().SendBuffMsg(MESSAGE_TAKINA, []byte("failed to unmarshal"))
		return
	}

	response := server(request)
	text, _ := json.Marshal(response)
	req.GetConnection().SendBuffMsg(MESSAGE_TAKINA, text)
}

func server(request TakinaRequest) TakinaResponse {
	if request.Token != global.Token {
		return TakinaResponse{
			Code: -1,
			Msg:  "token error",
		}
	}

	switch request.Type {
	case TAKINA_TYPE_ADD_PROXY:
		var data TakinaRequestStartProxy
		err := json.Unmarshal([]byte(request.Data), &data)
		if err != nil {
			return TakinaResponse{
				Code: -1,
				Msg:  "failed to unmarshal",
			}
		}
		raddr, rport, err := addProxy(data.Laddr, data.Lport, data.ProxyType)
		if err != nil {
			return TakinaResponse{
				Code: -1,
				Msg:  err.Error(),
			}
		}
		response := TakinaResponseStartProxy{
			Raddr: raddr,
			Rport: rport,
		}
		text, _ := json.Marshal(response)
		return TakinaResponse{
			Code: 0,
			Msg:  "success",
			Data: string(text),
		}
	case TAKINA_TYPE_DEL_PROXY:
		var data TakinaRequestStopProxy
		err := json.Unmarshal([]byte(request.Data), &data)
		if err != nil {
			return TakinaResponse{
				Code: -1,
				Msg:  "failed to unmarshal",
			}
		}
		err = delProxy(data.Laddr, data.Lport)
		if err != nil {
			return TakinaResponse{
				Code: -1,
				Msg:  err.Error(),
			}
		}
		response := TakinaResponseStopProxy{}
		text, _ := json.Marshal(response)
		return TakinaResponse{
			Code: 0,
			Msg:  "success",
			Data: string(text),
		}
	}

	return TakinaResponse{
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
