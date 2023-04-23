package server

import (
	"errors"
	"math/rand"

	"github.com/Yeuoly/Takina/src/frpcdaemon"
	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/proxy"
	"github.com/Yeuoly/Takina/src/types"
)

func (c *Takina) InitFrpcConfig(node TakinaNode) (*types.FrpcConfig, error) {
	frpc := &types.FrpcConfig{
		ServerAddr: node.Address,
	}

	resp, err := helper.SendGetAndParse[types.TakinaResponseWarp[types.TakinaResponseGetFrpsConfig]](
		node.GenerateUrl(types.ROUTER_TAKINA_SERVER_GET_FRPS_CONFIG),
		helper.HttpPayloadJson(GetPackedRequest(c, types.TakinaRequestGetFrpsConfig{})),
		helper.HttpTimeout(2000),
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
	c.Frpcs, err = frpcdaemon.LaunchFrpcDaemon(c.Frpcs, func() {
		// connect to frpc daemon
		proxy.LoadTakinaFrpc(c.Frpcs)
	})
	if err != nil {
		helper.Panic("[Takina] failed to launch frpc daemon: %s", err.Error())
	}

	helper.Info("[Takina] frpc daemon launched")
}

func (c *Takina) RandomNode() (*TakinaNode, error) {
	if len(c.Nodes) == 0 {
		return nil, errors.New("no node available")
	}

	idx := rand.Intn(len(c.Nodes))
	return &c.Nodes[idx], nil
}

func (c *Takina) GetNodeByAddress(address string) (*TakinaNode, error) {
	for _, node := range c.Nodes {
		if node.Address == address {
			return &node, nil
		}
	}

	return nil, errors.New("node not found")
}
