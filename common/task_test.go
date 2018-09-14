package common

import (
	"testing"
)

// Dummy ok backend
type DummyDriver struct{}

func (d *DummyDriver) IsExist(TaskUUID string) []string {
	return []string{}
}
func (d *DummyDriver) Kill(containerId string) {}
func (d *DummyDriver) Start(TaskUUID string, image string, cmd []string) string {
	return "01234"
}
func (d *DummyDriver) GetPort(containerID string) int {
	return 32123
}
func (d *DummyDriver) Clear() {}

// Error dummy backend
type DummyErrorDriver struct{}

func (d *DummyErrorDriver) IsExist(TaskUUID string) []string {
	return []string{}
}
func (d *DummyErrorDriver) Kill(containerId string) {}
func (d *DummyErrorDriver) Start(TaskUUID string, image string, cmd []string) string {
	return "01234"
}
func (d *DummyErrorDriver) GetPort(containerID string) int {
	return 32123
}
func (d *DummyErrorDriver) Clear() {}

// The test
func TestTask(t *testing.T) {
	// Valid state
	task := Task{
		ContainerDriver: &DummyDriver{},
		UUID:            "peisaipaishequaeroofeeyoghahwool",
	}
	task.AddContainer()
	task.DestroyAll()

	// Error state
	task = Task{
		ContainerDriver: &DummyErrorDriver{},
		UUID:            "peisaipaishequaeroofeeyoghahwool",
	}
	task.AddContainer()
	task.DestroyAll()
}
