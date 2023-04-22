package types

import (
	"encoding/json"
	"errors"
)

type TakinaRequest[T any] struct {
	Token string `json:"token"`
	Data  T      `json:"data"`
}

func ParseTakinaRequest[T any](data []byte) *TakinaRequest[T] {
	var request TakinaRequest[T]
	//unmarshal
	err := json.Unmarshal(data, &request)
	if err != nil {
		return nil
	}
	return &request
}

func (c *TakinaRequest[T]) Json() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func (c *TakinaRequest[T]) JsonBytes() []byte {
	data, _ := json.Marshal(c)
	return data
}

type TakinaResponseWarp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func (c *TakinaResponseWarp[T]) Success() bool {
	return c.Code == 0
}

func (c *TakinaResponseWarp[T]) Error() error {
	if c.Success() {
		return nil
	}
	return errors.New(c.Msg)
}

type TakinaResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (c *TakinaResponse) Json() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func (c *TakinaResponse) JsonBytes() []byte {
	data, _ := json.Marshal(c)
	return data
}

func SuccessResponse(data any) *TakinaResponse {
	return &TakinaResponse{Code: 0, Msg: "success", Data: data}
}

func ErrorResponse(code int, msg string) *TakinaResponse {
	return &TakinaResponse{Code: code, Msg: msg}
}
