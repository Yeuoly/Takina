package types

type InitTakinaDockerDaemonResponse struct {
	ContainerId   string `json:"container_id"`
	ContainerName string `json:"container_name"`
}
