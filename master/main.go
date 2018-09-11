package main

import (
	"github.com/labstack/echo"
)

var taskStorage TaskStorageInterface
var config MasterConfig

func init() {
	// Handle config
	ProcessEnvironmentVariables(&config)

	// Database connection
	taskStorage = &TaskStorage{
		DatabasePath: config.DatabaseFile,
	}
}

func main() {
	e := echo.New()
	e.POST("/tasks/", CreateTaskHandler)
	e.GET("/tasks/", ListTaskHandler)
	e.GET("/tasks/by-domain/:domain", GetTaskByDomainHandler)
	e.Logger.Fatal(e.Start(":1323"))
}
