package types

type TakinaRequestStartProxy struct {
	ProxyType string `json:"proxy_type"`
	Laddr     string `json:"laddr"`
	Lport     int    `json:"lport"`
}

type TakinaResponseStartProxy struct {
	Raddr string `json:"raddr"`
	Rport int    `json:"rport"`
}

type TakinaRequestStopProxy struct {
	Laddr string `json:"laddr"`
	Lport int    `json:"lport"`
}

type TakinaResponseStopProxy struct{}

type TakinaRequestListProxy struct{}

type TakinaResponseListProxy struct {
	Proxies []Proxy `json:"proxies"`
}
