package cli

import "github.com/Yeuoly/Takina/src/types"

func generateRequest[T any](data T) types.TakinaRequest[T] {
	return types.TakinaRequest[T]{
		Token: *takina_token,
		Data:  data,
	}
}
