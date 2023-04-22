package server

import (
	"github.com/Yeuoly/Takina/src/frpsdaemon"
	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
)

func (c *Takina) InitFrpsDaemonConfig() {
	c.Frps = &types.FrpsConfig{}
	c.Frps.BindAddr = "0.0.0.0"
	port, err := c.requestAvailablePort()
	if err != nil {
		helper.Panic("[Takina] failed to request available port: %s when launch frps", err.Error())
	}
	c.Frps.BindPort = port
}

func (c *Takina) RunFrpsDeamon() {
	c.InitFrpsDaemonConfig()

	var err error
	helper.Info("[Takina] launching frpc daemon...")
	c.Frps, err = frpsdaemon.LaunchFrpsDaemon(c.Frps)
	if err != nil {
		helper.Panic("[Takina] failed to launch frpc daemon: %s", err.Error())
	}
	helper.Info("[Takina] frpc daemon launched")
}

func (c *Takina) GetFrpsConfig() *types.FrpsConfig {
	return c.Frps
}
