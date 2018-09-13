package main

import (
	"github.com/labstack/echo"
	"strconv"
)

var taskStorage TaskStorageInterface
var config = MasterConfig{}

func main() {
	// Handle config
	ProcessEnvironmentVariables(&config)

	// Database connection
	taskStorage = &TaskStorage{
		DatabasePath: config.DatabaseFile,
	}
	db := taskStorage.open()
	defer db.Close()

	// TODO: Update tasks

	e := echo.New()
	e.POST("/tasks/", CreateTaskHandler)
	e.GET("/tasks/", ListTaskHandler)
	e.DELETE("/tasks/:uuid", DeleteTaskHandler)
	e.GET("/tasks/by-domain/:domain", GetTaskByDomainHandler)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Port)))
}
