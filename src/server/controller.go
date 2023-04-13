package server

import (
	"github.com/Yeuoly/Takina/src/types"
	"github.com/Yeuoly/zinx/ziface"
)

func BaiscController[T any](req ziface.IRequest, success func(T, ziface.IConnection)) {
	data := req.GetData()
	request := types.ParseTakinaRequest[T](data)
	if request == nil {
		req.GetConnection().Send(types.ErrorResponse("unsupported request").JsonBytes())
	} else {
		if !GetTakina().Auth(request.Token) {
			req.GetConnection().Send(types.ErrorResponse("invalid token").JsonBytes())
		} else {
			success(request.Data, req.GetConnection())
		}
	}
}

func (router *TakinaServerGetFrpsConfig) Handle(req ziface.IRequest) {
	BaiscController(req, func(data types.TakinaRequestGetFrpsConfig, conn ziface.IConnection) {
		config := GetTakina().GetFrpsConfig()
		conn.Send(types.SuccessResponse(types.TakinaResponseGetFrpsConfig{
			BindPort: config.BindPort,
			Token:    config.Token,
		}).JsonBytes())
	})
}

func (router *TakinaServerGetPort) Handle(req ziface.IRequest) {
	BaiscController(req, func(data types.TakinaRequestGetPort, conn ziface.IConnection) {
		port, err := GetTakina().requestAvailablePort()
		if err != nil {
			conn.Send(types.ErrorResponse(err.Error()).JsonBytes())
		} else {
			conn.Send(types.SuccessResponse(types.TakinaResponseGetPort{
				Port: port,
			}).JsonBytes())
		}
	})
}

func (couter *TakinaServerReleasePort) Handle(req ziface.IRequest) {
	BaiscController(req, func(data types.TakinaRequestReleasePort, conn ziface.IConnection) {
		err := GetTakina().releasePort(data.Port)
		if err != nil {
			conn.Send(types.ErrorResponse(err.Error()).JsonBytes())
		} else {
			conn.Send(types.SuccessResponse(types.TakinaResponseReleasePort{}).JsonBytes())
		}
	})
}
