package server

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Takina struct {
	Token            string `yaml:"token"`
	PortRange        string `yaml:"port_range"`
	ServerName       string `yaml:"server_name"`
	TakinaPort       int    `yaml:"takina_port"`
	RealPortRange    []int
	PortPool         *list.List
	requestPortMutex sync.Mutex
	Frps             *types.FrpsConfig
}

var global_takina_instance Takina

func init() {
	global_takina_instance.PortPool = list.New()
}

func InitTakinaServer() {
	data, err := ioutil.ReadFile("conf/takina_server.yaml")
	if err != nil {
		helper.Panic("[Takina] failed to read takina.yaml: %s", err.Error())
	}

	err = yaml.Unmarshal(data, &global_takina_instance)
	if err != nil {
		helper.Panic("[Takina] failed to unmarshal takina.yaml: %s", err.Error())
	}

	// parse port range
	parts := strings.Split(global_takina_instance.PortRange, ",")
	for _, part := range parts {
		ports := strings.Split(strings.TrimSpace(part), "-")
		if len(ports) == 1 {
			port, err := strconv.Atoi(ports[0])
			if err != nil {
				helper.Panic("[Takina] failed to parse port range: %s", err.Error())
			}
			global_takina_instance.RealPortRange = append(global_takina_instance.RealPortRange, port)
		} else if len(ports) == 2 {
			start, err := strconv.Atoi(ports[0])
			if err != nil {
				helper.Panic("[Takina] failed to parse port range: %s", err.Error())
			}
			end, err := strconv.Atoi(ports[1])
			if err != nil {
				helper.Panic("[Takina] failed to parse port range: %s", err.Error())
			}
			for i := start; i <= end; i++ {
				// check if port exists
				exists := false
				for _, port := range global_takina_instance.RealPortRange {
					if port == i {
						exists = true
						break
					}
				}
				if !exists {
					global_takina_instance.RealPortRange = append(global_takina_instance.RealPortRange, i)
				}
			}
		} else {
			helper.Panic("[Takina] failed to parse port range: %s", part)
		}
	}

	if len(global_takina_instance.RealPortRange) == 0 {
		helper.Panic("[Takina] Please ensure port range is correct")
	}

	// copy port range to port pool
	for _, port := range global_takina_instance.RealPortRange {
		global_takina_instance.PortPool.PushBack(port)
	}

	// check if server name is empty
	if global_takina_instance.ServerName == "" {
		helper.Panic("[Takina] Please ensure server name is correct")
	}

	// check if token is empty
	if global_takina_instance.Token == "" {
		helper.Panic("[Takina] Please ensure token is correct, empty token will cause security issue")
	}

	// check if takina port is empty
	if global_takina_instance.TakinaPort == 0 {
		helper.Panic("[Takina] Please ensure takina port is correct")
	}
}

func (root *Takina) setupRouter(r *gin.Engine) {
	r.GET(types.ROUTER_TAKINA_SERVER_GET_FRPS_CONFIG, TakinaServerGetFrpsConfig)
	r.POST(types.ROUTER_TAKINA_SERVER_REQUEST_PORT, TakinaServerGetPort)
	r.POST(types.ROUTER_TAKINA_SERVER_RELEASE_PORT, TakinaServerReleasePort)
}

func (root *Takina) Run() {
	// Launch frpc daemon
	root.RunFrpsDeamon()

	// launch server and listen
	server := gin.Default()
	root.setupRouter(server)
	gin.SetMode(gin.ReleaseMode)

	server.Run(fmt.Sprintf(":%d", root.TakinaPort))
}

func GetTakina() *Takina {
	return &global_takina_instance
}
