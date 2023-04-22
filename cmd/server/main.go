package main

import (
	server "github.com/Yeuoly/Takina/src/server"
)

func main() {
	server.InitTakinaServer()
	server.GetTakina().Run()
}
