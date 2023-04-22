package frpcdaemon

import (
	"fmt"
	"io"
	"sync"

	"github.com/Yeuoly/Takina/src/helper"
	"github.com/Yeuoly/Takina/src/types"
	go_memexec "github.com/amenzhinsky/go-memexec"
)

func generateFrpcConfig(config *types.FrpcConfig) string {
	return fmt.Sprintf(`[common]
server_addr = %s
server_port = %d
token = %s
admin_addr = %s
admin_port = %d
admin_user = %s
admin_pwd = %s

`, config.ServerAddr, config.ServerPort, config.Token, config.AdminAddr, config.AdminPort, config.AdminUser, config.AdminPwd)
}

// LaunchFrpcDaemon will launch all frpc daemon in the configs
// frpc daemon will be launched in serveral goroutines, so it will return immediately
// the return value is the configs with admin info
func LaunchFrpcDaemon(configs []*types.FrpcConfig) ([]*types.FrpcConfig, error) {
	for _, config := range configs {
		admin_user := helper.RandomStr(8)
		admin_pass := helper.RandomStr(16)
		admin_addr := "127.0.0.1"
		admin_port, err := helper.GetAvaliablePort()
		if err != nil {
			return nil, err
		}

		config.AdminAddr = admin_addr
		config.AdminPort = admin_port
		config.AdminUser = admin_user
		config.AdminPwd = admin_pass

		frpc_config_content := generateFrpcConfig(config)
		frpc_file_name, frpc_file, err := helper.CreateTempFile()
		if err != nil {
			return nil, err
		}

		n, err := frpc_file.Write([]byte(frpc_config_content))
		if err != nil || n != len(frpc_config_content) {
			frpc_file.Close()
			return nil, err
		}

		// close to flush to disk
		frpc_file.Close()

		// launch frpc
		frpc, err := go_memexec.New(frpcEmbed)
		if err != nil {
			return nil, err
		}

		cmd := frpc.Command("-c", frpc_file_name)
		go func() {
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				helper.Panic("[frpc] failed to get stdout: %s", err)
			}

			stderr, err := cmd.StderrPipe()
			if err != nil {
				helper.Panic("[frpc] failed to get stderr: %s", err)
			}

			err = cmd.Start()
			if err != nil {
				helper.Panic("[frpc] failed to start: %s", err)
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

			go read_pipe("frpc-log", "frpc", stdout)
			go read_pipe("frpc-error", "frpc", stderr)

			wg.Wait()

			err = cmd.Wait()
			if err != nil {
				helper.Warn("[frpc] failed to wait: %s", err)
			}
		}()
	}

	return configs, nil
}
