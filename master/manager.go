package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"shadowproject/common"
	shadowerrors "shadowproject/common/errors"
)

type TaskStorage struct {
	DatabasePath string
}

// Open the task storage
func (l *TaskStorage) open() *leveldb.DB {
	db, err := leveldb.OpenFile(l.DatabasePath, nil)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "cannot open the database",
		})
	}
	return db
}

// Add a new task
func (l *TaskStorage) Add(task *common.Task) {
	jsonBody, err := json.Marshal(task)
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "marshal error",
		})
	}

	db := l.open()
	defer db.Close()

	err = db.Put([]byte("task:"+task.UUID), jsonBody, &opt.WriteOptions{})
	if err != nil {
		panic(shadowerrors.ShadowError{
			Origin:         err,
			VisibleMessage: "database error",
		})
	}
}

// Return list of all tasks
func (l *TaskStorage) Filter() []common.Task {
	var tasks []common.Task
	var task common.Task

	db := l.open()
	defer db.Close()

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

// Filter one task based on a wanted domain
func (l *TaskStorage) GetByDomain(wantedDomain string) (*common.Task, error) {
	var task common.Task

	db := l.open()
	defer db.Close()

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
