package api

import _ "embed"

//go:embed client_cli
var clientCli []byte

//go:embed client_daemon
var clientDaemon []byte
