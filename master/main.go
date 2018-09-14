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

	// API endpoint
	e := echo.New()
	e.POST("/tasks/", CreateTaskHandler)
	e.GET("/tasks/", ListTaskHandler)
	e.PUT("/tasks/:uuid", UpdateTaskHandler)
	e.GET("/tasks/:uuid", GetTaskHandler)
	e.DELETE("/tasks/:uuid", DeleteTaskHandler)
	e.GET("/tasks/by-domain/:domain", GetTaskByDomainHandler)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Port)))
}
