package common

type ContainerDriver interface {
	Kill(containerId string) error
	Start(TaskUUID string) (string, error)
}
