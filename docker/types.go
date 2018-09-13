package docker

type ContainerDriverInterface interface {
	IsExist(TaskUUID string) []string
	Kill(containerId string)
	Start(TaskUUID string, image string, cmd []string) string
	GetPort(containerID string) int
	Clear()
}
