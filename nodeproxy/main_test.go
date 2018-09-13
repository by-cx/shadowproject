package main

import (
	"shadowproject/common"
)

//// Docker mock driver
type MockDockerDriver struct{}

func (m *MockDockerDriver) Kill(containerId string) error {
	return nil
}

func (m *MockDockerDriver) Start(TaskUUID string, image string, cmd []string) (string, error) {
	return "coomeeweebaibasaofiijengiefeejoh", nil
}

func (m *MockDockerDriver) GetPort(containerID string) (int, error) {
	return 32000, nil
}

func (m *MockDockerDriver) Clear() error {
	return nil
}

//// Mock client to test backend calls
type MockShadowMasterClient struct{}

func (m *MockShadowMasterClient) AddTask(domains []string, image string, command []string) (*common.Task, error) {
	return &common.Task{
		Driver:  &MockDockerDriver{},
		UUID:    "giajaiphobohroothoivaengukooquat",
		Domains: domains,
		Image:   image,
		Command: command,
	}, nil
}

func (m *MockShadowMasterClient) ListTasks() ([]common.Task, error) {
	return []common.Task{
		{
			Driver:  &MockDockerDriver{},
			UUID:    "giajaiphobohroothoivaengukooquat",
			Domains: []string{"localhost"},
			Image:   "shadow/testimage",
			Command: []string{"/srv/a_binary"},
		},
	}, nil
}

func (m *MockShadowMasterClient) GetTaskByDomain(wantedDomain string) (*common.Task, error) {
	return &common.Task{
		Driver:  &MockDockerDriver{},
		UUID:    "giajaiphobohroothoivaengukooquat",
		Domains: []string{"localhost"},
		Image:   "shadow/testimage",
		Command: []string{"/srv/a_binary"},
	}, nil
}

//// Set up the test environment
//func TestMain(m *testing.M) {
//	//shadowClient = &MockShadowMasterClient{}
//}
