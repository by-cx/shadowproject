package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"net/http"
	"shadowproject/common"
)

type ShadowMasterClient struct {
	Proto string
	Host  string
	Port  int
}

func (s *ShadowMasterClient) formatURL(action string) string {
	return fmt.Sprintf("%s://%s:%d%s", s.Proto, s.Host, s.Port, action)
}

func (s *ShadowMasterClient) AddTask(domains []string, image string, command []string) (*common.Task, error) {
	var task common.Task

	requestTask := &common.Task{
		Domains: domains,
		Image:   image,
		Command: command,
	}
	requestData, err := json.Marshal(requestTask)
	if err != nil {
		return &task, err
	}

	reader := bytes.NewReader(requestData)

	resp, err := grequests.Post(
		s.formatURL("/tasks/"),
		&grequests.RequestOptions{RequestBody: reader},
	)

	if err != nil {
		return &task, err
	}

	err = json.Unmarshal(resp.Bytes(), &task)
	if err != nil {
		panic(err)
	}

	return &task, nil
}

func (s *ShadowMasterClient) ListTasks() ([]common.Task, error) {
	var tasks []common.Task

	resp, err := grequests.Get(s.formatURL("/tasks/"), nil)
	if err != nil {
		return tasks, err
	}

	err = json.Unmarshal(resp.Bytes(), tasks)
	if err != nil {
		panic(err)
	}

	return tasks, nil
}

func (s *ShadowMasterClient) GetTask(TaskUUID string) (*common.Task, error) {
	var task common.Task

	resp, err := grequests.Get(s.formatURL("/tasks/"+TaskUUID), nil)
	if err != nil {
		return &task, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return &task, errors.New(resp.RawResponse.Status)
	}

	err = json.Unmarshal(resp.Bytes(), &task)
	if err != nil {
		panic(err)
	}

	return &task, nil
}

func (s *ShadowMasterClient) GetTaskByDomain(wantedDomain string) (*common.Task, error) {
	var task common.Task

	resp, err := grequests.Get(s.formatURL("/tasks/by-domain/"+wantedDomain), nil)
	if err != nil {
		return &task, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return &task, errors.New(resp.RawResponse.Status)
	}

	err = json.Unmarshal(resp.Bytes(), &task)
	if err != nil {
		panic(err)
	}

	return &task, nil
}
