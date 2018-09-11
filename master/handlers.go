package main

import (
	"github.com/labstack/echo"
	"net/http"
	"shadowproject/common"
)

// Handler to create a new task
func CreateTaskHandler(c echo.Context) error {
	var inputTask common.Task

	err := c.Bind(&inputTask)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, common.GeneralResponse{Message: err.Error()}, common.JSON_INDENT)
	}

	// Creating a new structure. The NewTask function validates the data by itself.
	task, errs := common.NewTask(inputTask.Domains, inputTask.Image, inputTask.Command)
	if len(errs) > 0 {
		message := ""
		for _, err := range errs {
			message += err.Error()
		}
		return c.JSONPretty(http.StatusBadRequest, common.GeneralResponse{Message: message}, common.JSON_INDENT)
	}

	taskStorage.Add(task)

	return c.JSONPretty(http.StatusOK, task, common.JSON_INDENT)
}

// Returns list of installed tasks
func ListTaskHandler(c echo.Context) error {
	tasks := taskStorage.Filter()
	return c.JSONPretty(http.StatusOK, tasks, common.JSON_INDENT)
}

// Handler to get complete task based on a domain
func GetTaskByDomainHandler(c echo.Context) error {
	wantedDomain := c.Param("domain")

	task, err := taskStorage.GetByDomain(wantedDomain)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, common.GeneralResponse{Message: err.Error()}, common.JSON_INDENT)
	}

	return c.JSONPretty(http.StatusOK, task, common.JSON_INDENT)
}
