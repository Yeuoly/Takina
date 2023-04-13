package main

import (
	client "github.com/Yeuoly/Takina/src/client"
)

func main() {
	client.InitTakinaClientDaemon()
	client.GetTakina().Run()
}