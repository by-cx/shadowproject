package common

type ContainerDriver interface {
	Kill(containerId string) error
	Start(TaskUUID string, image string, cmd []string) (string, error)
	GetPort(containerID string) (int, error)
}
