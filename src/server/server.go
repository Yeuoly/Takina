package server

import (
	"io/ioutil"

	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
	"github.com/Yeuoly/zinx/zlog"
	"github.com/Yeuoly/zinx/znet"
	"gopkg.in/yaml.v2"
)

type TakinaNode struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	Token   string `yaml:"token"`
}

type Takina struct {
	Token string       `yaml:"token"`
	Nodes []TakinaNode `yaml:"nodes"`
	Frpcs []*types.FrpcConfig
}

var global_takina_instance Takina

func InitTakinaClientDaemon() {
	data, err := ioutil.ReadFile("conf/takina.yaml")
	if err != nil {
		helper.Panic("[Takina] failed to read takina.yaml: %s", err.Error())
	}

	err = yaml.Unmarshal(data, &global_takina_instance)
	if err != nil {
		helper.Panic("[Takina] failed to unmarshal takina.yaml: %s", err.Error())
	}
}

func (root *Takina) Run() {
	// Launch frpc daemon
	root.RunFrpcDeamon()

	// launch zinx server and listen
	zlog.SetLogger(new(zinxLogger))
	server := znet.NewServer()
	server.AddRouter(MESSAGE_TAKINA, &TakinaServer{})
	server.Serve()
}

func GetTakina() *Takina {
	return &global_takina_instance
}
