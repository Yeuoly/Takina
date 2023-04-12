package main

import (
	"github.com/Yeuoly/Takina/src/server"
)

func main() {
	server.InitTakinaClientDaemon()
	server.GetTakina().Run()
}