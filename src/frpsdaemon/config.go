package frpsdaemon

import (
	"fmt"
	"io"
	"sync"

	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
	go_memexec "github.com/amenzhinsky/go-memexec"
)

func generateFrpsConfig(config *types.FrpsConfig) string {
	return fmt.Sprintf(`[common]
bind_port = %d
bind_addr = %s
token = %s
`, config.BindPort, config.BindAddr, config.Token)
}

// LaunchFrpsDaemon will launch frps daemon
// frps daemon will be launched in a goroutine, so it will return immediately
// the return value is the config with admin info
func LaunchFrpsDaemon(config *types.FrpsConfig) (*types.FrpsConfig, error) {
	token := helper.RandomStr(16)
	config.Token = token

	frpc_config_content := generateFrpsConfig(config)
	frps_file_name, frps_file, err := helper.CreateTempFile()
	if err != nil {
		return nil, err
	}

	n, err := frps_file.Write([]byte(frpc_config_content))
	if err != nil || n != len(frpc_config_content) {
		frps_file.Close()
		return nil, err
	}

	// close to flush to disk
	frps_file.Close()

	// launch frps
	frps, err := go_memexec.New(frpsEmbed)
	if err != nil {
		return nil, err
	}

	cmd := frps.Command("-c", frps_file_name)
	go func() {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			helper.Panic("[frps] failed to get stdout: %s", err)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			helper.Panic("[frps] failed to get stderr: %s", err)
		}

		err = cmd.Start()
		if err != nil {
			helper.Panic("[frps] failed to start: %s", err)
		}

		var wg sync.WaitGroup
		wg.Add(2)

		read_pipe := func(name string, module string, pipe io.ReadCloser) {
			defer wg.Done()
			// read stdout
			for {
				buf := make([]byte, 1024)
				n, err := pipe.Read(buf)
				if err != nil && err != io.EOF {
					helper.Warn("[%s][error] failed to read %s: %v", module, name, err)
					return
				}
				if n == 0 || err == io.EOF {
					break
				}
				helper.Info("[%s] %s", module, string(buf[:n]))
			}
		}

		go read_pipe("frps-log", "frpc", stdout)
		go read_pipe("frps-error", "frpc", stderr)

		wg.Wait()

		err = cmd.Wait()
		if err != nil {
			helper.Warn("[frps] failed to wait: %s", err)
		}
	}()

	return config, nil
}
