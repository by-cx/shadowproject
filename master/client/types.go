package client

import "shadowproject/common"

// Interface to cover master client library
type ShadowMasterClientInterface interface {
	AddTask(domains []string, image string, command []string) (*common.Task, error)
	ListTasks() ([]common.Task, error)
	GetTask(TaskUUID string) (*common.Task, error)
	GetTaskByDomain(wantedDomain string) (*common.Task, error)
}
