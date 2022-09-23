package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/aceld/zinx/znet"
)

type TakinaRequest struct {
	Token string `json:"token"`
	Type  string `json:"type"`
	Data  string `json:"data"`
}

type TakinaRequestStartProxy struct {
	ProxyType string `json:"proxy_type"`
	Laddr     string `json:"laddr"`
	Lport     int    `json:"lport"`
}

type TakinaRequestStopProxy struct {
	Laddr string `json:"laddr"`
	Lport int    `json:"lport"`
}

type TakinaResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type TakinaResponseStartProxy struct {
	Raddr string `json:"raddr"`
	Rport int    `json:"rport"`
}

type TakinaResponseStopProxy struct{}

const (
	TAKINA_TYPE_ADD_PROXY = "add_proxy"
	TAKINA_TYPE_DEL_PROXY = "del_proxy"
	TAKINA_TYPE_GET_PROXY = "get_proxy"
)

func main() {
	add()
}

func add() {
	addproxy_request := TakinaRequestStartProxy{
		ProxyType: "tcp",
		Laddr:     "127.0.0.1",
		Lport:     25570,
	}

	addproxy_request_json, _ := json.Marshal(addproxy_request)

	request := TakinaRequest{
		Token: "Takina",
		Type:  TAKINA_TYPE_ADD_PROXY,
		Data:  string(addproxy_request_json),
	}

	conn, err := net.Dial("tcp", "localhost:7171")
	if err != nil {
		panic(err)
	}

	close_chan := make(chan int, 1)

	go recv(conn, close_chan)

	text, _ := json.Marshal(request)
	dp := znet.NewDataPack()
	msg, _ := dp.Pack(znet.NewMsgPackage(1, text))
	_, err = conn.Write(msg)
	if err != nil {
		panic(err)
	}

	<-close_chan
	conn.Close()
}

func del() {
	delproxy_request := TakinaRequestStopProxy{
		Laddr: "127.0.0.1",
		Lport: 25570,
	}

	delproxy_request_json, _ := json.Marshal(delproxy_request)

	request := TakinaRequest{
		Token: "Takina",
		Type:  TAKINA_TYPE_DEL_PROXY,
		Data:  string(delproxy_request_json),
	}

	conn, err := net.Dial("tcp", "localhost:7171")
	if err != nil {
		panic(err)
	}

	close_chan := make(chan int, 1)

	go recv(conn, close_chan)

	text, _ := json.Marshal(request)
	dp := znet.NewDataPack()
	msg, _ := dp.Pack(znet.NewMsgPackage(1, text))
	_, err = conn.Write(msg)
	if err != nil {
		panic(err)
	}

	<-close_chan
	conn.Close()
}

func recv(conn net.Conn, close_chan chan int) {
	dp := znet.NewDataPack()

	head_data := make([]byte, dp.GetHeadLen())
	_, err := io.ReadFull(conn, head_data)
	if err != nil {
		return
	}

	msg_head, err := dp.Unpack(head_data)
	if err != nil {
		panic("failed to unpack message header")
	}

	if msg_head.GetDataLen() > 0 {
		msg := msg_head.(*znet.Message)
		msg.Data = make([]byte, msg.GetDataLen())

		_, err := io.ReadFull(conn, msg.Data)
		if err != nil {
			panic("failed to read message data")
		}

		fmt.Println(string(msg.Data))
	}

	close_chan <- 1
}
