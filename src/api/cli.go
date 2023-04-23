package api

import (
	"strconv"

	"github.com/Yeuoly/Takina/src/types"
)

func queryClient(args ...string) (string, error) {
	return dockerRunCommand("/client_cli", "", args...)
}

func StartProxy(laddr string, lport int, protocol string) (*types.TakinaClientCliStartProxyResponse, error) {
	result, err := queryClient(
		"-mode", "start",
		"-laddr", laddr,
		"-lport", strconv.Itoa(lport),
		"-type", protocol,
		"-takina_token", takina_token,
		"-fmt", "json",
	)

	if err != nil {
		return nil, err
	}

	response, err := parseJson[types.TakinaResponseWarp[types.TakinaClientCliStartProxyResponse]](
		result,
	)

	if err != nil {
		return nil, err
	}

	if response.Error() != nil {
		return nil, response.Error()
	}

	return &response.Data, nil
}

func StopProxy(laddr string, lport int) (*types.TakinaClientCliStopProxyResponse, error) {
	result, err := queryClient(
		"-mode", "stop",
		"-laddr", laddr,
		"-lport", strconv.Itoa(lport),
		"-takina_token", takina_token,
		"-fmt", "json",
	)

	if err != nil {
		return nil, err
	}

	response, err := parseJson[types.TakinaResponseWarp[types.TakinaClientCliStopProxyResponse]](
		result,
	)

	if err != nil {
		return nil, err
	}

	if response.Error() != nil {
		return nil, response.Error()
	}

	return &response.Data, nil
}

func ListProxy() (*types.TakinaClientCliListProxyResponse, error) {
	result, err := queryClient(
		"-mode", "list",
		"-takina_token", takina_token,
		"-fmt", "json",
	)

	if err != nil {
		return nil, err
	}

	response, err := parseJson[types.TakinaResponseWarp[types.TakinaClientCliListProxyResponse]](
		result,
	)

	if err != nil {
		return nil, err
	}

	if response.Error() != nil {
		return nil, response.Error()
	}

	return &response.Data, nil
}
