package types

type Proxy struct {
	Id    string `json:"id"`
	Laddr string `json:"laddr"`
	Lport int    `json:"lport"`
	Raddr string `json:"raddr"`
	Rport int    `json:"rport"`
	Type  string `json:"type"`
}
