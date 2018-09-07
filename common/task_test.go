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
func (d *DummyDriver) Start(TaskUUID string) (string, error) {
	return "01234", nil
}

// Error backend
type DummyErrorDriver struct{}

func (d *DummyErrorDriver) Kill(containerId string) error {
	return errors.New("Error!")
}
func (d *DummyErrorDriver) Start(TaskUUID string) (string, error) {
	return "01234", nil
}

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
