package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
	"math/rand"
	shadowerrors "shadowproject/common/errors"
	"strconv"
	"strings"
	"time"
)

const DOCKER_SOCK = "unix:///var/run/docker.sock"
const DOCKER_API_VERSION = "1.27"
const CONTAINER_DEFAULT_PORT = "8000/tcp"

type DockerDriver struct{}

func (d *DockerDriver) getClient() *dockerClient.Client {
	cli, err := dockerClient.NewClient(DOCKER_SOCK, DOCKER_API_VERSION, nil, nil)

	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "kill docker container error",
		})
	}

	return cli
}

func (d *DockerDriver) Kill(containerId string) {
	log.Println("Stopping container " + containerId)
	cli := d.getClient()

	timeout := time.Duration(30 * time.Second)
	err := cli.ContainerStop(context.TODO(), containerId, &timeout)

	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "kill docker container error",
		})
	}
}

// Checks existence of the container based on Task UUID.
// Returns container IDs in case of existence. Otherwise
// empty slice.
func (d *DockerDriver) IsExist(TaskUUID string) []string {
	var containerIDs = make([]string, 0)

	cli := d.getClient()

	containers, err := cli.ContainerList(context.TODO(), types.ContainerListOptions{})
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "docker error",
		})
	}
	for _, containerObject := range containers {
		for _, name := range containerObject.Names {
			name = strings.Trim(name, "/")
			if strings.Split(name, ".")[0] == TaskUUID {
				containerIDs = append(containerIDs, containerObject.ID)
			}
		}
	}

	return containerIDs
}

// Starts the container
func (d *DockerDriver) Start(TaskUUID string, image string, cmd []string) string {
	log.Println("Starting container " + TaskUUID)
	cli := d.getClient()

	portmaps := make(nat.PortMap, 1)
	portbindings := make([]nat.PortBinding, 1)
	portbindings[0] = nat.PortBinding{
		HostPort: "",
	}
	portmaps["8000/tcp"] = portbindings

	createdContainer, err := cli.ContainerCreate(
		context.TODO(),
		&container.Config{
			Hostname: TaskUUID,
			Env:      []string{},
			Image:    image,
			Cmd:      cmd,
		},
		&container.HostConfig{
			PortBindings: portmaps,
			AutoRemove:   true,
			Binds: []string{
				"/srv/" + TaskUUID + ":/srv",
			},
		},
		&network.NetworkingConfig{},
		TaskUUID+"."+strconv.Itoa(rand.Int()), // for multiple containers per task per server
	)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "start docker container error",
		})
	}

	containerId := createdContainer.ID

	err = cli.ContainerStart(context.TODO(), createdContainer.ID, types.ContainerStartOptions{})

	return containerId
}

// Return port we should redirect the request to
func (d *DockerDriver) GetPort(containerID string) int {
	cli := d.getClient()

	containerDetails, err := cli.ContainerInspect(context.TODO(), containerID)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "docker get port error",
		})
	}

	hostPort := containerDetails.NetworkSettings.Ports[CONTAINER_DEFAULT_PORT][0].HostPort
	port, err := strconv.Atoi(hostPort)

	return port
}

// Removes all containers the server contains
func (d *DockerDriver) Clear() {
	log.Println("Clearing docker containers")
	cli := d.getClient()

	timeout := time.Second * 30

	// Stop all containers
	containers, err := cli.ContainerList(context.TODO(), types.ContainerListOptions{})
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "clear docker containers error",
		})
	}
	for _, container := range containers {
		err = cli.ContainerStop(context.TODO(), container.ID, &timeout)
		if err != nil {
			panic(shadowerrors.ShadowError{
				Origin:         err,
				VisibleMessage: "clear docker containers error",
			})
		}
	}

	// Remove all containers
	containers, err = cli.ContainerList(context.TODO(), types.ContainerListOptions{})
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "clear docker containers error",
		})
	}
	for _, container := range containers {
		err = cli.ContainerRemove(context.TODO(), container.ID, types.ContainerRemoveOptions{})
		if err != nil {
			panic(shadowerrors.ShadowError{
				Origin:         err,
				VisibleMessage: "clear docker containers error",
			})
		}
	}
}
