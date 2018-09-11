package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"shadowproject/common"
)

// Interface to cover storage for tasks
type TaskStorageInterface interface {
	open() *leveldb.DB
	Add(task *common.Task)
	Filter() []common.Task
	GetByDomain(wantedDomain string) (*common.Task, error)
}
