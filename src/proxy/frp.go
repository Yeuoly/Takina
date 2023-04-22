package proxy

import (
	"strconv"
	"strings"
	"sync"

	"github.com/Yeuoly/Takina/src/types"
)

type Proxy = types.Proxy

func GenerateConfigContent(c *Proxy) string {
	content := `

[` + c.Id + `]
type = ` + c.Type + `
local_ip = ` + c.Laddr + `
local_port = ` + strconv.Itoa(c.Lport) + `
remote_port = ` + strconv.Itoa(c.Rport) + `

`
	return content
}

type FrpcNote struct {
	Uuid           string `yaml:"uuid"`
	Address        string `yaml:"address"`
	RAddress       string `yaml:"raddress"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Pass           string `yaml:"pass"`
	OriginalConfig string
	CurrentProxy   map[string]Proxy
	mtx            sync.RWMutex
}

type FrpcTcpConfig struct {
	Err        string `json:"err"`
	LocalAddr  string `json:"local_addr"`
	Name       string `json:"name"`
	RemoteAddr string `json:"remote_addr"`
	Plugin     string `json:"plugin"`
	Status     string `json:"status"`
	Type       string `json:"type"`
}

func (c *FrpcTcpConfig) parseAddress(address string) (string, int) {
	if address == "" {
		return "", 0
	}
	port_text := address[strings.LastIndex(address, ":")+1:]
	port, err := strconv.Atoi(port_text)
	if err != nil {
		return "", 0
	}
	return address[:strings.LastIndex(address, ":")], port
}

func (c *FrpcTcpConfig) LocalPort() int {
	_, port := c.parseAddress(c.LocalAddr)
	return port
}

func (c *FrpcTcpConfig) RemotePort() int {
	_, port := c.parseAddress(c.RemoteAddr)
	return port
}

func (c *FrpcTcpConfig) LocalAddress() string {
	address, _ := c.parseAddress(c.LocalAddr)
	return address
}

func (c *FrpcTcpConfig) RemoteAddress() string {
	address, _ := c.parseAddress(c.RemoteAddr)
	return address
}

type TakinaConfig struct {
	ServerName    string `yaml:"server-name"`
	ServerAddress string `yaml:"server-address"`
	ServerPort    int    `yaml:"server-port"`
	AuthKey       string `yaml:"auth-key"`

	//frpc admin
	ClientNotes []FrpcNote `yaml:"client-notes"`
}

var globalConfig TakinaConfig

func loadDefaultProxy() {
	// load orginal config and current proxy
	for i := range globalConfig.ClientNotes {
		globalConfig.ClientNotes[i].CurrentProxy = make(map[string]Proxy)

		config_content := GetFrpConfig(&globalConfig.ClientNotes[i])
		if config_content == "" {
			continue
		}
		//load only [common] section
		config_content = config_content[strings.Index(config_content, "[common]")+8:]

		//fetch server address from config as remote address
		server_address := config_content[strings.Index(config_content, "server_addr = ")+14:]
		server_address = server_address[:strings.Index(server_address, "\n")]
		globalConfig.ClientNotes[i].RAddress = server_address

		if strings.Contains(config_content, "[") {
			config_content = config_content[:strings.Index(config_content, "[")]
		}

		//ensure only 1 line break in the end
		if strings.HasSuffix(config_content, "\n\n") {
			config_content = config_content[:len(config_content)-1]
		}

		globalConfig.ClientNotes[i].OriginalConfig = "[common]" + config_content
		result := GetFrpStatus(&globalConfig.ClientNotes[i])
		for _, v := range result.Tcp {
			globalConfig.ClientNotes[i].CurrentProxy[v.Name] = Proxy{
				Id:    v.Name,
				Laddr: v.LocalAddress(),
				Lport: v.LocalPort(),
				Raddr: v.RemoteAddress(),
				Rport: v.RemotePort(),
				Type:  v.Type,
			}
		}
	}
}

func GetConfig() TakinaConfig {
	return globalConfig
}
