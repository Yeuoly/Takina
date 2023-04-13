package server

func (c *Takina) Auth(token string) bool {
	return c.Token == token
}
