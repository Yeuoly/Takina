package main

import (
	"fmt"

	"github.com/Yeuoly/Takina/src/proxy"
)

func main() {
	//err := proxy.StopProxy("127.0.0.1", 25561)
	raddr, rpotr, err := proxy.AutoLaunchProxy("127.0.0.1", 25570, "tcp")
	fmt.Println(raddr, rpotr, err)
	fmt.Println(proxy.GetProxies())
}
