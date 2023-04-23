package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yeuoly/Takina/src/api"
)

func main() {
	config := `server-name: Takina # server name for takina
token: testtest # token use to conmunicate between cli and daemon-client
takina_port: 40002 # takina port

nodes: 
  - 
    address: 172.17.0.1 # takina server address
    port: 40001 # takina server port
    token: testtest # takina server token
`
	r, err := api.InitTakinaDockerDaemon("testtest", strings.NewReader(config), func(s string) {
		//fmt.Printf("[TakinaDockerDaemon] %s", s)
	}, func(s string) {
		//fmt.Printf("[TakinaDockerDaemon] %s", s)
	})

	fmt.Println(r, err)

	// wait for daemon to start
	time.Sleep(5 * time.Second)

	resp, err := api.StartProxy("127.0.0.1", 80, "tcp")
	fmt.Println(resp, err)
	list, err := api.ListProxy()
	fmt.Println(list, err)
	resp1, err := api.StopProxy("127.0.0.1", 80)
	fmt.Println(resp1, err)
	list, err = api.ListProxy()
	fmt.Println(list, err)

	var x string
	fmt.Scanf("%s", &x)
	api.StopTakinaDockerDaemon()
}
