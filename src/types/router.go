package types

const (
	ROUTER_TAKINA_CLIENT_DAEMON_START_PROXY = iota
	ROUTER_TAKINA_CLIENT_DAEMON_STOP_PROXY
	ROUTER_TAKINA_CLIENT_DAEMON_LIST_PROXY

	ROUTER_TAKINA_SERVER_GET_FRPS_CONFIG = iota
	ROUTER_TAKINA_SERVER_REQUEST_PORT
	ROUTER_TAKINA_SERVER_RELEASE_PORT
)