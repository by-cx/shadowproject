package common

import (
	"github.com/pkg/errors"
	"testing"
)

// Dummy ok backend
type DummyDriver struct{}

func (d *DummyDriver) Kill(containerId string) error {
	return nil
}
func (d *DummyDriver) Start(TaskUUID string, image string, cmd []string) (string, error) {
	return "01234", nil
}

func (d *DummyDriver) GetPort(containerID string) (int, error) {
	return 32123, nil
}

func (d *DummyDriver) Clear() error {
	return nil
}

// Error dummy backend
type DummyErrorDriver struct{}

func (d *DummyErrorDriver) Kill(containerId string) error {
	return errors.New("Error!")
}
func (d *DummyErrorDriver) Start(TaskUUID string, image string, cmd []string) (string, error) {
	return "01234", nil
}

func (d *DummyErrorDriver) GetPort(containerID string) (int, error) {
	return 32123, nil
}

func (d *DummyErrorDriver) Clear() error {
	return nil
}

// The test
func TestTask(t *testing.T) {
	// Valid state
	task := Task{
		Driver: &DummyDriver{},
		UUID:   "peisaipaishequaeroofeeyoghahwool",
	}
	task.AddContainer()
	task.DestroyAll()

	// Error state
	task = Task{
		Driver: &DummyErrorDriver{},
		UUID:   "peisaipaishequaeroofeeyoghahwool",
	}
	task.AddContainer()
	task.DestroyAll()
}
