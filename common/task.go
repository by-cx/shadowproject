package common

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"shadowproject/docker"
	"strings"
)

type Task struct {
	Driver docker.ContainerDriverInterface `json:"-"` // Container driver for managing Containers

	UUID       string   `json:"uuid"`    // Identification of the task
	Containers []string `json:"-"`       // Container's docker IDs
	Domains    []string `json:"domains"` // Domain list on which this tasks listens
	Image      string   `json:"image"`   // Docker image
	Command    []string `json:"command"` // Command to run
}

func NewTask(domains []string, image string, command []string) (*Task, []error) {
	var task Task

	taskUUID := uuid.NewV4()
	task.Driver = &docker.DockerDriver{}
	task.UUID = strings.Replace(taskUUID.String(), "-", "", -1)
	task.Domains = domains
	task.Image = image
	task.Command = command

	errorList := task.Validate()

	return &task, errorList
}

// Task data validation
func (t *Task) Validate() []error {
	var errorList []error

	if len(t.UUID) != 32 {
		errorList = append(errorList, errors.New("the UUID has to be 32 characters long string"))
	}
	if len(t.Image) == 0 {
		errorList = append(errorList, errors.New("image can't be empty"))
	}
	if len(t.Command) == 0 {
		errorList = append(errorList, errors.New("command can't be empty"))
	}
	if len(t.Domains) == 0 {
		errorList = append(errorList, errors.New("define at least one domain"))
	}

	return errorList
}

// Adds new container for this task. Returns container ID and error.
func (t *Task) AddContainer() (string, error) {
	containerId, err := t.Driver.Start(t.UUID, t.Image, t.Command)
	if err == nil {
		t.Containers = append(t.Containers, containerId)
	}
	return containerId, err
}

func (t *Task) DestroyAll() error {
	var remainingIds []string
	var lastErr error

	for _, containerId := range t.Containers {
		err := t.Driver.Kill(containerId)
		if err != nil {
			lastErr = err
			remainingIds = append(remainingIds, containerId)
		}
	}

	t.Containers = remainingIds

	return lastErr
}
