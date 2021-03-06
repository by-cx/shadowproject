package main

import (
	"shadowproject/common"
)

//// Docker mock driver
type MockDockerDriver struct{}

func (m *MockDockerDriver) Status() bool {
	return true
}
func (m *MockDockerDriver) IsExist(TaskUUID string) []string {
	return []string{}
}

func (m *MockDockerDriver) Kill(containerId string) {}

func (m *MockDockerDriver) Start(TaskUUID string, image string, cmd []string, target string) string {
	return "coomeeweebaibasaofiijengiefeejoh"
}

func (m *MockDockerDriver) GetPort(containerID string) int {
	return 80
}

func (m *MockDockerDriver) Clear() {}

//// Mock client to test backend calls
type MockShadowMasterClient struct {
	ReturnedErr error
}

func (m *MockShadowMasterClient) AddTask(domains []string, image string, command []string) (*common.Task, error) {
	return &common.Task{
		ContainerDriver: &MockDockerDriver{},
		UUID:            "giajaiphobohroothoivaengukooquat",
		Domains:         domains,
		Image:           image,
		Command:         command,
		VolumeType:      "S3",
		Source:          "test/code.zip",
	}, m.ReturnedErr
}

func (m *MockShadowMasterClient) ListTasks() ([]common.Task, error) {
	return []common.Task{
		{
			ContainerDriver: &MockDockerDriver{},
			UUID:            "giajaiphobohroothoivaengukooquat",
			Domains:         []string{"localhost"},
			Image:           "shadow/testimage",
			Command:         []string{"/srv/a_binary"},
			VolumeType:      "S3",
			Source:          "test/code.zip",
		},
	}, m.ReturnedErr
}

func (m *MockShadowMasterClient) GetTask(taskUUID string) (*common.Task, error) {
	return &common.Task{
		ContainerDriver: &MockDockerDriver{},
		UUID:            "giajaiphobohroothoivaengukooquat",
		Domains:         []string{"localhost"},
		Image:           "shadow/testimage",
		Command:         []string{"/srv/a_binary"},
		VolumeType:      "S3",
		Source:          "test/code.zip",
	}, m.ReturnedErr
}

func (m *MockShadowMasterClient) GetTaskByDomain(wantedDomain string) (*common.Task, error) {
	return &common.Task{
		ContainerDriver: &MockDockerDriver{},
		UUID:            "giajaiphobohroothoivaengukooquat",
		Domains:         []string{"localhost"},
		Image:           "shadow/testimage",
		Command:         []string{"/srv/a_binary"},
		VolumeType:      "S3",
		Source:          "test/code.zip",
	}, m.ReturnedErr
}

// S3 volume driver mock
type MockS3VolumeDriver struct{}

func (m *MockS3VolumeDriver) Mount(source string, target string) {}
func (m *MockS3VolumeDriver) Umount(target string)               {}
