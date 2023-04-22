package server

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type TakinaNode struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	Token   string `yaml:"token"`
}

func (c *TakinaNode) GenerateUrl(path string) string {
	return fmt.Sprintf("http://%s:%d%s", c.Address, c.Port, path)
}

type Takina struct {
	Token      string       `yaml:"token"`
	Nodes      []TakinaNode `yaml:"nodes"`
	TakinaPort int          `yaml:"takina_port"`
	Frpcs      []*types.FrpcConfig
}

var global_takina_instance Takina

func InitTakinaClientDaemon() {
	data, err := ioutil.ReadFile("conf/takina_client.yaml")
	if err != nil {
		helper.Panic("[Takina] failed to read takina.yaml: %s", err.Error())
	}

	err = yaml.Unmarshal(data, &global_takina_instance)
	if err != nil {
		helper.Panic("[Takina] failed to unmarshal takina.yaml: %s", err.Error())
	}

	if global_takina_instance.TakinaPort == 0 {
		helper.Panic("[Takina] takina_port is not set")
	}
}

func (root *Takina) setupRouter(server *gin.Engine) {
	server.POST(types.ROUTER_TAKINA_CLIENT_DAEMON_START_PROXY, TakinaClientDeamonRequestStartProxy)
	server.POST(types.ROUTER_TAKINA_CLIENT_DAEMON_STOP_PROXY, TakinaClientDeamonRequestStopProxy)
	server.GET(types.ROUTER_TAKINA_CLIENT_DAEMON_LIST_PROXY, TakinaClientDeamonRequestListProxy)
}

func (root *Takina) Run() {
	// Launch frpc daemon
	root.RunFrpcDeamon()

	// launch server and listen
	server := gin.Default()
	root.setupRouter(server)
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	var err error
	gin.DefaultWriter, err = os.Create(os.DevNull)
	if err != nil {
		helper.Panic("[Takina] Failed to ignore log info: " + err.Error())
	}

	server.Run(fmt.Sprintf(":%d", root.TakinaPort))
}

func GetTakina() *Takina {
	return &global_takina_instance
}
