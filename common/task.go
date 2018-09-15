package common

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"shadowproject/common/containers"
	shadowerrors "shadowproject/common/errors"
	"shadowproject/common/volumes"
	"strings"
	"time"
)

const VolumeTypeS3 = "S3"

type Task struct {
	ContainerDriver containers.ContainerDriverInterface `json:"-"` // Container driver for managing Containers
	VolumeDriver    volumes.VolumeInterface             `json:"-"` // Volume driver to manage volumes

	UUID       string   `json:"uuid"`        // Identification of the task
	LastUpdate int64    `json:"last_update"` // Timestamp of the last change`
	Domains    []string `json:"domains"`     // Domain list on which this tasks listens
	Image      string   `json:"image"`       // Docker image
	Command    []string `json:"command"`     // Command to run
	VolumeType string   `json:"volume_type"` // Volume type: possible choices: S3
	Source     string   `json:"source"`      // Where the source code is located, for example path to file in S3 bucket
}

func NewTask(domains []string, image string, command []string, volumeType string, source string) (*Task, []error) {
	var task Task

	taskUUID := uuid.NewV4()
	task.ContainerDriver = &containers.DockerDriver{}
	task.UUID = strings.Replace(taskUUID.String(), "-", "", -1)
	task.Domains = domains
	task.Image = image
	task.Command = command
	task.Source = source
	task.VolumeType = volumeType
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
	if t.VolumeType != VolumeTypeS3 {
		errorList = append(errorList, errors.New("this volume type is no available, posible choices: "+VolumeTypeS3))
	}

	return errorList
}

// Adds new container for this task. Returns container ID and error.
func (t *Task) AddContainer() string {
	target := "/srv/" + t.UUID
	err := os.MkdirAll(target, 0755)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "mounting error",
		})
	}
	t.VolumeDriver.Mount(t.Source, target)

	containerId := t.ContainerDriver.Start(t.UUID, t.Image, t.Command, target)
	return containerId
}

// Remove all containers related to this task
func (t *Task) DestroyAll() {
	target := "/srv/" + t.UUID

	log.Println(t.ContainerDriver.IsExist(t.UUID))
	for _, containerId := range t.ContainerDriver.IsExist(t.UUID) {
		log.Println("Debug: killing", containerId, "created for", t.UUID)
		// Kill the container
		t.ContainerDriver.Kill(containerId)
		// Umount the volume
		t.VolumeDriver.Umount(target)
	}
}
