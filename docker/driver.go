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
	"strconv"
	"time"
)

const DOCKER_SOCK = "unix:///var/run/docker.sock"
const DOCKER_API_VERSION = "1.27"
const CONTAINER_DEFAULT_PORT = "8000/tcp"

type DockerDriver struct {
}

func (d *DockerDriver) getClient() (*dockerClient.Client, error) {
	cli, err := dockerClient.NewClient(DOCKER_SOCK, DOCKER_API_VERSION, nil, nil)
	return cli, err
}

func (d *DockerDriver) Kill(containerId string) error {
	log.Println("Stoping container " + containerId)
	cli, err := d.getClient()
	if err != nil {
		return err
	}

	timeout := time.Duration(30 * time.Second)
	err = cli.ContainerStop(context.TODO(), containerId, &timeout)

	return err
}

func (d *DockerDriver) Start(TaskUUID string, image string, cmd []string) (string, error) {
	log.Println("Starting container " + TaskUUID)
	cli, err := d.getClient()
	if err != nil {
		return "", err
	}

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
		TaskUUID+strconv.Itoa(rand.Int()),
	)
	if err != nil {
		return "", err
	}

	containerId := createdContainer.ID

	err = cli.ContainerStart(context.TODO(), createdContainer.ID, types.ContainerStartOptions{})

	return containerId, err
}

func (d *DockerDriver) GetPort(containerID string) (int, error) {
	cli, err := d.getClient()
	if err != nil {
		return 0, err
	}

	containerDetails, err := cli.ContainerInspect(context.TODO(), containerID)
	if err != nil {
		return 0, err
	}

	hostPort := containerDetails.NetworkSettings.Ports[CONTAINER_DEFAULT_PORT][0].HostPort
	port, err := strconv.Atoi(hostPort)

	return port, err
}
