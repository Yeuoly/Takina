package types

type FrpcConfig struct {
	// frpc server address
	ServerAddr string `json:"server_addr" yaml:"server_addr"`
	// frpc server port
	ServerPort int `json:"server_port" yaml:"server_port"`
	// frpc token
	Token string `json:"token" yaml:"token"`
	// frpc admin address
	AdminAddr string `json:"admin_addr" yaml:"admin_addr"`
	// frpc admin port
	AdminPort int `json:"admin_port" yaml:"admin_port"`
	// frpc admin user
	AdminUser string `json:"admin_user" yaml:"admin_user"`
	// frpc admin password
	AdminPwd string `json:"admin_pwd" yaml:"admin_pwd"`
}

var (
	frpcs []*FrpcConfig
)

func GetFrpcConfig() []*FrpcConfig {
	return frpcs
}

func SetFrpcConfig(configs []*FrpcConfig) {
	frpcs = configs
}

func AddFrpcConfig(config *FrpcConfig) {
	frpcs = append(frpcs, config)
}

func DelFrpcConfig(index int) {
	if index < 0 || index >= len(frpcs) {
		return
	}
	frpcs = append(frpcs[:index], frpcs[index+1:]...)
}
