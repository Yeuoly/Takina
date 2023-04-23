package api

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func dockerRunCommand(elf_path string, stdin string, args ...string) (result string, err error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return "", err
	}
	//run command in container
	args = append([]string{elf_path}, args...)
	exec_config := types.ExecConfig{
		Cmd:          args,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		User:         "root",
	}

	resp, err := cli.ContainerExecCreate(context.Background(), takina_container_id, exec_config)
	if err != nil {
		return "", err
	}

	hijack, err := cli.ContainerExecAttach(context.Background(), resp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	//set timeout
	hijack.Conn.SetDeadline(time.Now().Add(time.Minute))

	//write stdin
	go func() {
		defer hijack.CloseWrite()
		hijack.Conn.Write([]byte(stdin))
	}()

	//read
	result = ""
	last := make([]byte, 0)
	for {
		//check if last is already contain a complete output
		if len(last) >= 8 {
			//read length
			length := binary.BigEndian.Uint32(last[4:8])
			if len(last) >= 8+int(length) {
				result += string(last[8 : 8+length])
				last = last[8+length:]
				continue
			}
		}

		buf := make([]byte, 1024)
		n, err := hijack.Conn.Read(buf)
		if err != nil {
			break
		}
		//append buf to last
		if len(last) == 0 {
			last = buf[:n]
		} else {
			last = append(last, buf[:n]...)
		}
	}
	return
}
