package types

type FrpsConfig struct {
	BindAddr string `json:"bind_addr"`
	BindPort int    `json:"bind_port"`
	// frps token
	Token string `json:"token"`
}
