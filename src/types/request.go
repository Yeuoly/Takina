package types

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