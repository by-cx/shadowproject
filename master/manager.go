package main

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
	"shadowproject/common"
)

type TaskStorage struct {
	DatabasePath string
}

func (l *TaskStorage) open() *leveldb.DB {
	db, err := leveldb.OpenFile(l.DatabasePath, nil)
	if err != nil {
		log.Panicln(err)
	}
	return db
}

func (l *TaskStorage) Add(task *common.Task) {
	jsonBody, err := json.Marshal(task)
	if err != nil {
		log.Panicln(err)
	}

	db := l.open()
	defer db.Close()

	err = db.Put([]byte("task:"+task.UUID), jsonBody, &opt.WriteOptions{})
	if err != nil {
		log.Panicln(err)
	}
}

func (l *TaskStorage) Filter() ([]common.Task, error) {
	var tasks []common.Task
	var task common.Task

	db := l.open()
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		err := json.Unmarshal(iter.Value(), &task)
		if err != nil {
			log.Panicln(err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
