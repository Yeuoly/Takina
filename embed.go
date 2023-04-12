package takina

import _ "embed"

//go:embed embed/frpc
var FrpcEmbed []byte

//go:embed embed/frps
var FrpsEmbed []byte
