package helper

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Yeuoly/zinx/ziface"
	"github.com/Yeuoly/zinx/znet"
)

/*
	Zinx request helper
*/

type zinxRequestOptions struct {
	Type  string
	Value any
}

func ZinxRequestWithTimeout(timeout time.Duration) *zinxRequestOptions {
	return &zinxRequestOptions{
		Type:  "timeout",
		Value: timeout,
	}
}

type zinxResponseRouter struct {
	znet.BaseRouter
	Handler func(ziface.IRequest)
}

func (r *zinxResponseRouter) Handle(request ziface.IRequest) {
	r.Handler(request)
}

// One-shot request, send request and parse response like http
func SendZinxAndParse[R any](addr string, port int, route int, request any, options ...*zinxRequestOptions) (*R, error) {
	timeout := time.Duration(time.Second * 5)
	for _, option := range options {
		switch option.Type {
		case "timeout":
			timeout = option.Value.(time.Duration)
		}
	}

	timer := time.NewTimer(timeout)
	finish_chan := make(chan struct{}, 1)
	var resposne R
	var resposne_err error

	onZinxClientConnected := func(conn ziface.IConnection) {
		// json encode request
		request_text, err := json.Marshal(request)
		if err != nil {
			resposne_err = errors.New("failed to encode request")
			finish_chan <- struct{}{}
			return
		}

		conn.SendBuffMsg(uint32(route), request_text)
	}

	client := znet.NewClient(addr, port)
	client.SetOnConnStart(onZinxClientConnected)
	client.AddRouter(uint32(route), &zinxResponseRouter{
		Handler: func(request ziface.IRequest) {
			// get response
			response_text := request.GetData()
			err := json.Unmarshal(response_text, &resposne)
			if err != nil {
				resposne_err = errors.New("failed to decode response")
				finish_chan <- struct{}{}
				return
			}
			finish_chan <- struct{}{}
		},
	})
	client.Start()
	defer client.Stop()

	select {
	case <-timer.C:
		return nil, errors.New("timeout")
	case <-finish_chan:
		return &resposne, resposne_err
	}
}
