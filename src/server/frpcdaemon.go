package server

import (
	"github.com/Yeuoly/Takina/src/frpcdaemon"
	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
)

func (c *Takina) InitFrpcDaemonConfig() {
	for _, node := range c.Nodes {
		frpc := types.FrpcConfig{
			ServerAddr: node.Address,
			ServerPort: node.Port,
			Token:      node.Token,
		}
		c.Frpcs = append(c.Frpcs, &frpc)
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
