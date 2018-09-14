package common

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"log"
	"shadowproject/docker"
	"strings"
	"time"
)

type Task struct {
	Driver docker.ContainerDriverInterface `json:"-"` // Container driver for managing Containers

	UUID       string   `json:"uuid"`        // Identification of the task
	LastUpdate int64    `json:"last_update"` // Timestamp of the last change`
	Domains    []string `json:"domains"`     // Domain list on which this tasks listens
	Image      string   `json:"image"`       // Docker image
	Command    []string `json:"command"`     // Command to run
	Source     string   `json:"source"`      // Where the source code is located on our S3 bucket
}

func NewTask(domains []string, image string, command []string, source string) (*Task, []error) {
	var task Task

	taskUUID := uuid.NewV4()
	task.Driver = &docker.DockerDriver{}
	task.UUID = strings.Replace(taskUUID.String(), "-", "", -1)
	task.Domains = domains
	task.Image = image
	task.Command = command
	task.Source = source
	task.LastUpdate = time.Now().Unix()

	errorList := task.Validate()

	return &task, errorList
}

func (t *Task) Update(domains []string, image string, command []string, source string) []error {
	t.Domains = domains
	t.Image = image
	t.Command = command
	t.Source = source
	t.LastUpdate = time.Now().Unix()

	errorList := t.Validate()

	return errorList
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
	if len(t.Source) == 0 {
		errorList = append(errorList, errors.New("path to the source has to be defined"))
	}

	return errorList
}

// Adds new container for this task. Returns container ID and error.
func (t *Task) AddContainer() string {
	containerId := t.Driver.Start(t.UUID, t.Image, t.Command)
	return containerId
}

func (t *Task) DestroyAll() {
	log.Println(t.Driver.IsExist(t.UUID))
	for _, containerId := range t.Driver.IsExist(t.UUID) {
		log.Println("Debug: killing", containerId, "created for", t.UUID)
		t.Driver.Kill(containerId)
	}
}
