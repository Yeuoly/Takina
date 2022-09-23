package proxy

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Yeuoly/Takina/src/helper"
)

type Api struct{}

type StatusResponse struct {
	Http  interface{}     `json:"http"`
	Https interface{}     `json:"https"`
	Tcp   []FrpcTcpConfig `json:"tcp"`
	Udp   interface{}     `json:"udp"`
	Stcp  interface{}     `json:"stcp"`
	Xtcp  interface{}     `json:"xtcp"`
	Sudp  interface{}     `json:"sudp"`
}

func (a *Api) basicHeader(user string, pwd string) map[string]string {
	return map[string]string{
		"Authorization": "Basic " + helper.Base64Encode(user+":"+pwd),
	}
}

func (a *Api) parseResult(src string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(src), &result)
	if err != nil {
		fmt.Errorf("parse result error: %s", err)
	}
	return result
}

func (a *Api) getFrpStatus(note *FrpcNote) string {
	resp, err := helper.HttpGet("http://"+note.Address+":"+strconv.Itoa(note.Port)+"/api/status", a.basicHeader(note.User, note.Pass))
	if err != nil {
		fmt.Errorf("get frp status error: %s", err)
		return ""
	}

	return resp
}

func (a *Api) getFrpConfig(note *FrpcNote) string {
	resp, err := helper.HttpGet("http://"+note.Address+":"+strconv.Itoa(note.Port)+"/api/config", a.basicHeader(note.User, note.Pass))
	if err != nil {
		fmt.Errorf("get frp config error: %s", err)
		return ""
	}
	return resp
}

func (a *Api) putFrpConfig(note *FrpcNote, config string) (string, error) {
	resp, err := helper.HttpPut("http://"+note.Address+":"+strconv.Itoa(note.Port)+"/api/config", config, a.basicHeader(note.User, note.Pass))
	if err != nil {
		fmt.Errorf("put frp config error: %s", err)
		return "", err
	}
	return resp, nil
}

func (a *Api) reloadFrp(note *FrpcNote) (string, error) {
	resp, err := helper.HttpGet("http://"+note.Address+":"+strconv.Itoa(note.Port)+"/api/reload", a.basicHeader(note.User, note.Pass))
	if err != nil {
		fmt.Errorf("reload frp error: %s", err)
		return "", err
	}
	return resp, nil
}

var defaultApi = &Api{}

func GetFrpStatus(note *FrpcNote) StatusResponse {
	var status StatusResponse
	result := defaultApi.getFrpStatus(note)
	json.Unmarshal([]byte(result), &status)
	return status
}

func GetFrpConfig(note *FrpcNote) string {
	return defaultApi.getFrpConfig(note)
}

func PutFrpConfig(note *FrpcNote, config string) (string, error) {
	return defaultApi.putFrpConfig(note, config)
}

func ReloadFrp(note *FrpcNote) (string, error) {
	return defaultApi.reloadFrp(note)
}
