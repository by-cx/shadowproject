package docker

type ContainerDriverInterface interface {
	Kill(containerId string) error
	Start(TaskUUID string, image string, cmd []string) (string, error)
	GetPort(containerID string) (int, error)
	Clear() error
}
