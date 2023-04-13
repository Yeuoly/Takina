package server

import (
	"github.com/Yeuoly/Takina/src/frpsdaemon"
	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
)

func (c *Takina) InitFrpsDaemonConfig() {
	c.Frps.BindAddr = "0.0.0.0"
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