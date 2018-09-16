package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Dummy ok backend
type DummyDriver struct{}

func (d *DummyDriver) IsExist(TaskUUID string) []string {
	return []string{}
}
func (d *DummyDriver) Kill(containerId string) {}
func (d *DummyDriver) Start(TaskUUID string, image string, cmd []string, target string) string {
	return "01234"
}
func (d *DummyDriver) GetPort(containerID string) int {
	return 32123
}
func (d *DummyDriver) Clear() {}

// Mock volume driver
type MockVolumeDriver struct{}

func (m *MockVolumeDriver) Mount(source string, target string) {}
func (m *MockVolumeDriver) Umount(target string)               {}

// The test
func TestTask(t *testing.T) {
	// Valid state
	task := Task{
		ContainerDriver: &DummyDriver{},
		VolumeDriver:    &MockVolumeDriver{},
		UUID:            "peisaipaishequaeroofeeyoghahwool",
		Image:           "testimage",
		Source:          "source/code.zip",
		Domains:         []string{"localhost:1234"},
		VolumeType:      "S3",
		Command:         []string{"/srv/testbin"},
	}
	task.AddContainer()
	task.DestroyAll()

	// Validation
	task.UUID = "peisaipaishequaeroofeeyoghahwool additional something"
	firstError := task.Validate()[0]
	assert.Equal(t, "the UUID has to be 32 characters long string", firstError.Error())
	task.UUID = "peisaipaishequaeroofeeyoghahwool"

	task.Image = ""
	firstError = task.Validate()[0]
	assert.Equal(t, "image can't be empty", firstError.Error())
	task.Image = "testimage"

	task.Source = ""
	firstError = task.Validate()[0]
	assert.Equal(t, "path to the source has to be defined", firstError.Error())
	task.Source = "source/code.zip"

	task.Domains = []string{}
	firstError = task.Validate()[0]
	assert.Equal(t, "define at least one domain", firstError.Error())
	task.Domains = []string{"localhost:1234"}

	task.VolumeType = ""
	firstError = task.Validate()[0]
	assert.Equal(t, "this volume type is no available, posible choices: S3", firstError.Error())
	task.VolumeType = "S3"

	task.Command = []string{}
	firstError = task.Validate()[0]
	assert.Equal(t, "command can't be empty", firstError.Error())
	task.Command = []string{"/srv/testbin"}
}

func TestNewTask(t *testing.T) {
	task, errorList := NewTask([]string{"localhost:1234"}, "testimage", []string{"/srv/testbin"}, "S3", "test/code.zip")

	assert.Equal(t, []string{"localhost:1234"}, task.Domains)
	assert.Equal(t, "testimage", task.Image)
	assert.Equal(t, []string{"/srv/testbin"}, task.Command)
	assert.Equal(t, "S3", task.VolumeType)
	assert.Equal(t, "test/code.zip", task.Source)
	assert.NotEqual(t, 0, task.LastUpdate)
	assert.Equal(t, 32, len(task.UUID))
	assert.Equal(t, 0, len(errorList))
}

func TestTask_Update(t *testing.T) {
	task, _ := NewTask([]string{"localhost:1234"}, "testimage", []string{"/srv/testbin"}, "S3", "test/code.zip")

	assert.Equal(t, []string{"localhost:1234"}, task.Domains)
	assert.Equal(t, "testimage", task.Image)
	assert.Equal(t, []string{"/srv/testbin"}, task.Command)
	assert.Equal(t, "S3", task.VolumeType)
	assert.Equal(t, "test/code.zip", task.Source)

	task.Update([]string{"localhost:1235"}, "testimage2", []string{"/srv/testbin2"}, "test/code2.zip")

	assert.Equal(t, []string{"localhost:1235"}, task.Domains)
	assert.Equal(t, "testimage2", task.Image)
	assert.Equal(t, []string{"/srv/testbin2"}, task.Command)
	assert.Equal(t, "test/code2.zip", task.Source)
}
