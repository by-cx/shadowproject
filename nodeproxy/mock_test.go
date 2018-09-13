package main

import (
	"shadowproject/common"
)

//// Docker mock driver
type MockDockerDriver struct {
	ReturnedErr error
}

func (m *MockDockerDriver) Kill(containerId string) error {
	return m.ReturnedErr
}

func (m *MockDockerDriver) Start(TaskUUID string, image string, cmd []string) (string, error) {
	return "coomeeweebaibasaofiijengiefeejoh", m.ReturnedErr
}

func (m *MockDockerDriver) GetPort(containerID string) (int, error) {
	return 32000, m.ReturnedErr
}

func (m *MockDockerDriver) Clear() error {
	return m.ReturnedErr
}

//// Mock client to test backend calls
type MockShadowMasterClient struct {
	ReturnedErr error
}

func (m *MockShadowMasterClient) AddTask(domains []string, image string, command []string) (*common.Task, error) {
	return &common.Task{
		Driver:  &MockDockerDriver{},
		UUID:    "giajaiphobohroothoivaengukooquat",
		Domains: domains,
		Image:   image,
		Command: command,
	}, m.ReturnedErr
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
	}, m.ReturnedErr
}

func (m *MockShadowMasterClient) GetTaskByDomain(wantedDomain string) (*common.Task, error) {
	return &common.Task{
		Driver:  &MockDockerDriver{},
		UUID:    "giajaiphobohroothoivaengukooquat",
		Domains: []string{"localhost"},
		Image:   "shadow/testimage",
		Command: []string{"/srv/a_binary"},
	}, m.ReturnedErr
}
