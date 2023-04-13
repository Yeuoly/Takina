package types

type TakinaRequestGetFrpsConfig struct{}

type TakinaResponseGetFrpsConfig struct {
	BindPort int    `json:"bind_port"`
	Token    string `json:"token"`
}

type TakinaRequestGetPort struct{}

type TakinaResponseGetPort struct {
	Port int `json:"port"`
}

type TakinaRequestReleasePort struct {
	Port int `json:"port"`
}

type TakinaResponseReleasePort struct{}
