package server

import "github.com/Yeuoly/Takina/src/types"

func (c *Takina) Auth(token string) bool {
	return c.Token == token
}

func NewRequest[T any](c *Takina, data T) *types.TakinaRequest[T] {
	return &types.TakinaRequest[T]{
		Token: c.Token,
		Data:  data,
	}
}
