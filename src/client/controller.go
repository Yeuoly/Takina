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

func (router *TakinaClientDeamonRequestStartProxy) Handle(req ziface.IRequest) {
	BaiscController(req, func(data types.TakinaRequestStartProxy, conn ziface.IConnection) {
		raddr, rport, err := addProxy(data.Laddr, data.Lport, data.ProxyType)
		if err != nil {
			conn.Send(types.ErrorResponse(err.Error()).JsonBytes())
		} else {
			conn.Send(types.SuccessResponse(types.TakinaResponseStartProxy{
				Raddr: raddr,
				Rport: rport,
			}).JsonBytes())
		}
	})
}

func (router *TakinaClientDeamonRequestStopProxy) Handle(req ziface.IRequest) {
	BaiscController(req, func(data types.TakinaRequestStopProxy, conn ziface.IConnection) {
		err := delProxy(data.Laddr, data.Lport)
		if err != nil {
			conn.Send(types.ErrorResponse(err.Error()).JsonBytes())
		} else {
			conn.Send(types.SuccessResponse(types.TakinaResponseStopProxy{}).JsonBytes())
		}
	})
}

func (router *TakinaClientDeamonRequestListProxy) Handle(req ziface.IRequest) {
	BaiscController(req, func(data types.TakinaRequestListProxy, conn ziface.IConnection) {
		proxies := listProxy()
		conn.Send(types.SuccessResponse(types.TakinaResponseListProxy{
			Proxies: proxies,
		}).JsonBytes())
	})
}
