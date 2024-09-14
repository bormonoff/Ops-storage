package storage

import (
	"ops-storage/internal/server/storage/memory"
	"ops-storage/internal/server/storage/postgres"
)

var store StorageFacade

type StorageFacade interface {
	Insert(counterType string, name string, val string) error
	Get(counterType string, name string) (string, error)
	GetAll() (*map[string]string, error)

	IsStorageAlive() bool
}

type RecoverConfig struct {
	RelPath  string
	Interval int
	Restore  bool
}

func Init(dsn string, cfg RecoverConfig) {
	if dsn == "" {
		memory.InitRecover(cfg.RelPath, cfg.Interval, cfg.Restore)
		store = memory.Instance()
	} else {
		store = postgres.New(dsn)
	}
}

func Instance() StorageFacade {
	return store
}
