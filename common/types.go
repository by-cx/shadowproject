package common

const JSON_INDENT = "  "

type ContainerDriver interface {
	Kill(containerId string) error
	Start(TaskUUID string, image string, cmd []string) (string, error)
	GetPort(containerID string) (int, error)
	Clear() error
}

type GeneralResponse struct {
	Message string `json:"message"`
}
