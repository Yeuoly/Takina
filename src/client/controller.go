package server

import (
	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
	"github.com/gin-gonic/gin"
)

func BaiscController[T any](r *gin.Context, success func(T)) {
	helper.BindRequest(r, func(request types.TakinaRequest[T]) {
		if !GetTakina().Auth(request.Token) {
			r.JSON(200, types.ErrorResponse(-403, "token error"))
		} else {
			success(request.Data)
		}
	})
}

func TakinaClientDeamonRequestStartProxy(r *gin.Context) {
	BaiscController(r, func(data types.TakinaRequestStartProxy) {
		raddr, rport, err := addProxy(data.Laddr, data.Lport, data.ProxyType)
		if err != nil {
			r.JSON(200, types.ErrorResponse(-500, err.Error()))
		} else {
			r.JSON(200, types.SuccessResponse(types.TakinaResponseStartProxy{
				Raddr: raddr,
				Rport: rport,
			}))
		}
	})
}

func TakinaClientDeamonRequestStopProxy(r *gin.Context) {
	BaiscController(r, func(data types.TakinaRequestStopProxy) {
		err := delProxy(data.Laddr, data.Lport)
		if err != nil {
			r.JSON(200, types.ErrorResponse(-500, err.Error()))
		} else {
			r.JSON(200, types.SuccessResponse(types.TakinaResponseStopProxy{}))
		}
	})
}

func TakinaClientDeamonRequestListProxy(r *gin.Context) {
	BaiscController(r, func(data types.TakinaRequestListProxy) {
		proxies := listProxy()
		r.JSON(200, types.SuccessResponse(types.TakinaResponseListProxy{
			Proxies: proxies,
		}))
	})
}
