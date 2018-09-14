package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"shadowproject/common"
	shadowerrors "shadowproject/common/errors"
	"strings"
)

type TaskStorage struct {
	DatabasePath string
	db           *leveldb.DB
}

// TODO: this needs refactoring. Get misses not found feature and error handling should be similar to other parts f the code.

// Open the task storage
func (l *TaskStorage) open() *leveldb.DB {
	if l.db == nil {
		db, err := leveldb.OpenFile(l.DatabasePath, nil)
		if err != nil {
			panic(shadowerrors.ShadowError{
				Origin:         err,
				VisibleMessage: "cannot open the database",
			})
		}
		l.db = db
		return db
	}

	return l.db
}

// Add a new task
func (l *TaskStorage) Add(task *common.Task) error {
	jsonBody, err := json.Marshal(task)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "marshal error",
		})
	}

	db := l.open()

	var duplicatedDomains []string
	for _, domain := range task.Domains {
		if l.CheckDomainDuplicity(domain) != "" {
			duplicatedDomains = append(duplicatedDomains, domain)
		}
	}
	if len(duplicatedDomains) != 0 {
		return errors.New("these domains are already in the system: " + strings.Join(duplicatedDomains, ", "))
	}

	err = db.Put([]byte("task:"+task.UUID), jsonBody, &opt.WriteOptions{})
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "add to database error",
		})
	}

	return nil
}

// Update parameters of the task with taskUUID
func (l *TaskStorage) Update(taskUUID string, domains []string, image string, command []string, source string) (*common.Task, []error) {
	db := l.open()

	// Find the task
	task, err := l.Get(taskUUID)
	if err != nil {
		return nil, []error{err}
	}

	// Try to update the task
	errorList := task.Update(domains, image, command, source)

	if len(errorList) == 0 {
		// Check if domains are duplicated
		var duplicatedDomains []string
		for _, domain := range task.Domains {
			if l.CheckDomainDuplicity(domain) != "" {
				duplicatedDomains = append(duplicatedDomains, domain)
			}
		}
		if len(duplicatedDomains) != 0 {
			return task, []error{errors.New("these domains are already in the system: " + strings.Join(duplicatedDomains, ", "))}
		}

		// Make a JSON from it
		jsonBody, err := json.Marshal(task)
		if err != nil {
			panic(shadowerrors.ShadowError{
				Origin:         err,
				VisibleMessage: "marshal error",
			})
		}

		// Replace the old version of the task
		err = db.Put([]byte("task:"+task.UUID), jsonBody, &opt.WriteOptions{})
		if err != nil {
			panic(shadowerrors.ShadowError{
				Origin:         err,
				VisibleMessage: "update in database error",
			})
		}
	}

	return task, errorList
}

// Return list of all tasks
func (l *TaskStorage) Filter() []common.Task {
	var tasks []common.Task
	var task common.Task

	db := l.open()

	iter := db.NewIterator(util.BytesPrefix([]byte("task:")), nil)
	for iter.Next() {
		err := json.Unmarshal(iter.Value(), &task)
		if err != nil {
			panic(shadowerrors.ShadowError{
				Origin:         err,
				VisibleMessage: "marshal error",
			})
		}
		tasks = append(tasks, task)
	}

	return tasks
}

// Delete a task
func (l *TaskStorage) Delete(taskUUID string) error {
	db := l.open()

	return db.Delete([]byte("task:"+taskUUID), &opt.WriteOptions{})
}

// Returns one task by its UUID
func (l *TaskStorage) Get(TaskUUID string) (*common.Task, error) {
	var task common.Task

	db := l.open()
	value, err := db.Get([]byte("task:"+TaskUUID), &opt.ReadOptions{})

	err = json.Unmarshal(value, &task)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "marshal error",
		})
	}

	// TODO: implement not found error
	return &task, err
}

// Filter one task based on a wanted domain
func (l *TaskStorage) GetByDomain(wantedDomain string) (*common.Task, error) {
	var task common.Task

	db := l.open()

	iter := db.NewIterator(util.BytesPrefix([]byte("task:")), nil)
	for iter.Next() {
		err := json.Unmarshal(iter.Value(), &task)
		if err != nil {
			panic(shadowerrors.ShadowError{
				Origin:         err,
				VisibleMessage: "marshal error",
			})
		}

		for _, domain := range task.Domains {
			if domain == wantedDomain {
				return &task, nil
			}
		}
	}

	return nil, errors.New("no task found")
}

// Check for domain duplicity
// Returns task id in case it finds the given domain in the db
// Empty string means not found
func (l *TaskStorage) CheckDomainDuplicity(domainToCheck string) string {
	var task common.Task

	db := l.open()

	iter := db.NewIterator(util.BytesPrefix([]byte("task:")), nil)
	for iter.Next() {
		err := json.Unmarshal(iter.Value(), &task)
		if err != nil {
			panic(shadowerrors.ShadowError{
				Origin:         err,
				VisibleMessage: "marshal error",
			})
		}

		for _, domain := range task.Domains {
			if domain == domainToCheck {
				return strings.Split(string(iter.Key()), ":")[1]
			}
		}
	}

	return ""
}
