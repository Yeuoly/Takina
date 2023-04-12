package server

import (
	"encoding/json"

	"github.com/Yeuoly/Takina/src/types"
	"github.com/Yeuoly/zinx/ziface"
	"github.com/Yeuoly/zinx/znet"
)

const MESSAGE_TAKINA = iota

type TakinaServer struct {
	znet.BaseRouter
}

func (router *TakinaServer) Handle(req ziface.IRequest) {
	var request types.TakinaRequest
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
