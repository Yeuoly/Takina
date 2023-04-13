package server

import (
	"time"

	"github.com/Yeuoly/Takina/src/frpcdaemon"
	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
)

func (c *Takina) InitFrpcConfig(node TakinaNode) (*types.FrpcConfig, error) {
	frpc := &types.FrpcConfig{
		ServerAddr: node.Address,
	}

	resp, err := helper.SendZinxAndParse[types.TakinaResponseWarp[types.TakinaResponseGetFrpsConfig]](
		node.Address,
		node.Port,
		types.ROUTER_TAKINA_SERVER_GET_FRPS_CONFIG,
		helper.ZinxRequestWithTimeout(time.Second*2),
	)
	if err != nil {
		return nil, err
	}

	err = resp.Error()
	if err != nil {
		return nil, err
	}

	frpc.ServerPort = resp.Data.BindPort
	frpc.Token = resp.Data.Token

	return frpc, nil
}

func (c *Takina) InitFrpcDaemonConfig() {
	for _, node := range c.Nodes {
		frpc, err := c.InitFrpcConfig(node)
		if err != nil {
			helper.Panic("[Takina] failed to init frpc config: %s", err.Error())
		}

		c.Frpcs = append(c.Frpcs, frpc)
	}
}

func (c *Takina) RunFrpcDeamon() {
	c.InitFrpcDaemonConfig()

	var err error
	helper.Info("[Takina] launching frpc daemon...")
	c.Frpcs, err = frpcdaemon.LaunchFrpcDaemon(c.Frpcs)
	if err != nil {
		helper.Panic("[Takina] failed to launch frpc daemon: %s", err.Error())
	}
	helper.Info("[Takina] frpc daemon launched")
}