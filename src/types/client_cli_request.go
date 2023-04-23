package types

type TakinaClientCliStartProxyResponse struct {
	Proxy Proxy `json:"proxy"`
}

type TakinaClientCliStopProxyResponse struct{}

type TakinaClientCliListProxyResponse struct {
	Proxies []Proxy `json:"proxies"`
}
