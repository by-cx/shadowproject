package main

import (
	"github.com/labstack/echo"
	"net/http"
	"shadowproject/common"
)

var taskStorage *TaskStorage

func init() {
	taskStorage = &TaskStorage{
		DatabasePath: "tasks.db",
	}
}

// Handler to create a new task
func CreateTaskHandler(c echo.Context) error {
	var inputTask common.Task

	err := c.Bind(&inputTask)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, common.GeneralResponse{Message: err.Error()}, "  ")
	}

	// Creating a new structure. The NewTask function validates the data by itself.
	task, errs := common.NewTask(inputTask.Domains, inputTask.Image, inputTask.Command)
	if len(errs) > 0 {
		message := ""
		for _, err := range errs {
			message += err.Error()
		}
		return c.JSONPretty(http.StatusBadRequest, common.GeneralResponse{Message: message}, "  ")
	}

	taskStorage.Add(task)

	return c.JSONPretty(http.StatusOK, task, "  ")
}

// Returns list of installed tasks
func ListTaskHandler(c echo.Context) error {
	tasks, err := taskStorage.Filter()
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, common.GeneralResponse{Message: err.Error()}, "  ")
	}

	return c.JSONPretty(http.StatusOK, tasks, "  ")
}

func main() {
	e := echo.New()
	e.POST("/tasks/", CreateTaskHandler)
	e.GET("/tasks/", ListTaskHandler)
	e.Logger.Fatal(e.Start(":1323"))
}
