package server

import (
	"container/list"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
	"github.com/Yeuoly/zinx/zlog"
	"github.com/Yeuoly/zinx/znet"
	"gopkg.in/yaml.v2"
)

type Takina struct {
	Token            string `yaml:"token"`
	PortRange        string `yaml:"port_range"`
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
}

func (root *Takina) Run() {
	// Launch frpc daemon
	root.RunFrpsDeamon()

	// launch zinx server and listen
	zlog.SetLogger(new(zinxLogger))
	server := znet.NewServer()
	server.Serve()
}

func GetTakina() *Takina {
	return &global_takina_instance
}
