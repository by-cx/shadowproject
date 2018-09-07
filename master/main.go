package main

import (
	"shadowproject/common"
	"shadowproject/docker"
	"time"
)

func main() {
	task := common.Task{
		Driver: &docker.DockerDriver{},
		UUID:   "peisaipaishequaeroofeeyoghahwool",
	}
	task.AddContainer()
	task.AddContainer()
	task.AddContainer()

	time.Sleep(time.Second * 20)

	task.DestroyAll()
}
