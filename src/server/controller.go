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

func TakinaServerGetFrpsConfig(r *gin.Context) {
	BaiscController(r, func(request types.TakinaRequestGetFrpsConfig) {
		config := GetTakina().GetFrpsConfig()
		r.JSON(200, types.SuccessResponse(types.TakinaResponseGetFrpsConfig{
			BindPort: config.BindPort,
			Token:    config.Token,
		}))
	})
}

func TakinaServerGetPort(r *gin.Context) {
	BaiscController(r, func(request types.TakinaRequestGetPort) {
		port, err := GetTakina().requestAvailablePort()
		if err != nil {
			r.JSON(200, types.ErrorResponse(-500, err.Error()))
		} else {
			r.JSON(200, types.SuccessResponse(types.TakinaResponseGetPort{
				Port: port,
			}))
		}
	})
}

func TakinaServerReleasePort(r *gin.Context) {
	BaiscController(r, func(request types.TakinaRequestReleasePort) {
		err := GetTakina().releasePort(request.Port)
		if err != nil {
			r.JSON(200, types.ErrorResponse(-500, err.Error()))
		} else {
			r.JSON(200, types.SuccessResponse(types.TakinaResponseReleasePort{}))
		}
	})
}
