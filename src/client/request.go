package server

import "github.com/Yeuoly/Takina/src/types"

func GetPackedRequest[T any](root *Takina, data T) types.TakinaRequest[T] {
	return types.TakinaRequest[T]{
		Token: root.Token,
		Data:  data,
	}
}
