package api

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path"

	"github.com/Yeuoly/Takina/src/helper"
	takina_types "github.com/Yeuoly/Takina/src/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

/*
 * Takina API provides methods to manage Takina Client.
 */
var takina_token string
var takina_container_id string

// InitTakinaDockerDaemon initializes Takina Docker Daemon, which is a Docker container running in the background.
// And it will launch a goroutine to manage TakinaClientDaemon.
func InitTakinaDockerDaemon(token string, config_file io.Reader, message_callback func(string), error_callback func(string)) (*takina_types.InitTakinaDockerDaemonResponse, error) {
	takina_token = token

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return nil, err
	}

	// check if alpine exists
	_, _, err = cli.ImageInspectWithRaw(ctx, "alpine")
	if err != nil {
		// pull alpine image
		resp, err := cli.ImagePull(ctx, "alpine", types.ImagePullOptions{})
		if err != nil {
			return nil, err
		}

		// read response
		for {
			buffer := make([]byte, 1024)
			n, err := resp.Read(buffer)
			if err != nil && err != io.EOF {
				error_callback(err.Error())
			}
			if err != nil {
				error_callback(err.Error())
			}
			if n == 0 || err == io.EOF {
				break
			}
			message_callback(string(buffer))
		}
	}

	init_configuration := true

	temp_file_path, remover, err := helper.CreateTempDir()
	if err != nil {
		return nil, err
	}
	defer remover()

	// write to temp file
	client_cli_reader := bytes.NewReader(clientCli)
	client_daemon_reader := bytes.NewReader(clientDaemon)
	config_file_path := path.Join(temp_file_path, "conf")
	helper.WriteToFile(config_file_path, "takina_client.yaml", config_file)
	helper.WriteToFile(temp_file_path, "client_cli", client_cli_reader)
	helper.WriteToFile(temp_file_path, "client_daemon", client_daemon_reader)

	client_cli_path := path.Join(temp_file_path, "client_cli")
	client_daemon_path := path.Join(temp_file_path, "client_daemon")
	config_file_path = path.Join(temp_file_path, "conf", "takina_client.yaml")

	// from file to tar
	tar_file_path, tar_file, err := helper.CreateTempFile()
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tar_file_path)

	// create a gzip writer
	gzip_writer, err := gzip.NewWriterLevel(tar_file, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	defer gzip_writer.Close()

	// create tar writer
	tar_writer := tar.NewWriter(gzip_writer)
	defer tar_writer.Close()

	// read file
	client_cli_file, err := os.Open(client_cli_path)
	if err != nil {
		return nil, err
	}
	defer client_cli_file.Close()

	client_daemon_file, err := os.Open(client_daemon_path)
	if err != nil {
		return nil, err
	}
	defer client_daemon_file.Close()

	client_config_file, err := os.Open(config_file_path)
	if err != nil {
		return nil, err
	}
	defer client_config_file.Close()

	tar_writer_func := func(file *os.File, filename string) error {
		//write file to tar
		stat, err := file.Stat()
		if err != nil {
			return err
		}
		client_cli_file_header, err := tar.FileInfoHeader(stat, filename)
		if err != nil {
			return err
		}
		client_cli_file_header.Name = filename
		if err := tar_writer.WriteHeader(client_cli_file_header); err != nil {
			return err
		}
		if _, err := io.Copy(tar_writer, file); err != nil {
			return err
		}
		return nil
	}

	if err := tar_writer_func(client_cli_file, "client_cli"); err != nil {
		return nil, err
	}
	if err := tar_writer_func(client_daemon_file, "client_daemon"); err != nil {
		return nil, err
	}
	if err := tar_writer_func(client_config_file, "conf/takina_client.yaml"); err != nil {
		return nil, err
	}

	tar_writer.Close()
	gzip_writer.Close()
	tar_file.Close()

	tar_file_context, err := os.Open(tar_file_path)
	if err != nil {
		return nil, err
	}
	defer tar_file_context.Close()

	// launch container
	c, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "alpine",
			Cmd:   []string{"tail", "-f", "/dev/null"},
		},
		&container.HostConfig{
			DNS:        []string{"8.8.8.8", "114.114.114.114"},
			Init:       &init_configuration,
			AutoRemove: true,
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		"takina",
	)

	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	// copy client daemon binary to container
	if err := cli.CopyToContainer(ctx, c.ID, "/", tar_file_context, types.CopyToContainerOptions{}); err != nil {
		return nil, err
	}

	// chmod client binary
	exec, err := cli.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
		Cmd: []string{"chmod", "+x", "/client_cli"},
	})
	if err != nil {
		return nil, err
	}
	_, err = cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}

	// chmod client daemon binary
	exec, err = cli.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
		Cmd: []string{"chmod", "+x", "/client_daemon"},
	})
	if err != nil {
		error_callback(err.Error())
		return nil, err
	}
	_, err = cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}

	// mkdir /conf
	exec, err = cli.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
		Cmd: []string{"mkdir", "/conf"},
	})
	if err != nil {
		error_callback(err.Error())
		return nil, err
	}
	_, err = cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}

	// launch client daemon
	go func() {
		exec, err := cli.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
			Cmd:          []string{"/client_daemon"},
			User:         "root",
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
		})
		if err != nil {
			error_callback(err.Error())
			return
		}

		hijack, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
		if err != nil {
			error_callback(err.Error())
			return
		}

		last := make([]byte, 0)
		for {
			//check if last is already contain a complete output
			if len(last) >= 8 {
				//read length
				length := binary.BigEndian.Uint32(last[4:8])
				if len(last) >= 8+int(length) {
					message_callback(string(last[8 : 8+length]))
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
	}()

	takina_container_id = c.ID

	return &takina_types.InitTakinaDockerDaemonResponse{
		ContainerId:   c.ID,
		ContainerName: "takina",
	}, nil
}

func StopTakinaDockerDaemon() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return err
	}

	// find container named takina
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	container_id := ""
	for _, container := range containers {
		if container.Names[0] == "/takina" {
			container_id = container.ID
		}
	}
	if container_id == "" {
		return errors.New("container not found")
	}
	if err := cli.ContainerStop(ctx, container_id, container.StopOptions{}); err != nil {
		return err
	}

	return nil
}
