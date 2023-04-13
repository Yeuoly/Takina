package server

import (
	"github.com/Yeuoly/zinx/znet"
)

type TakinaClientDeamonRequestStartProxy struct {
	znet.BaseRouter
}

type TakinaClientDeamonRequestStopProxy struct {
	znet.BaseRouter
}

type TakinaClientDeamonRequestListProxy struct {
	znet.BaseRouter
}
