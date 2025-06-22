package persistence

import (
	"github.com/boss-net/api/cloud/oss"

	"github.com/boss-net/api/boss-plugin/internal/types/app"
	"github.com/boss-net/api/boss-plugin/internal/utils/log"
)

var (
	persistence *Persistence
)

func InitPersistence(oss oss.OSS, config *app.Config) {
	persistence = &Persistence{
		storage:        NewWrapper(oss, config.PersistenceStoragePath),
		maxStorageSize: config.PersistenceStorageMaxSize,
	}

	log.Info("Persistence initialized")
}

func GetPersistence() *Persistence {
	return persistence
}
