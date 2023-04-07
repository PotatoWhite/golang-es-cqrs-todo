package query

import (
	"github.com/potato/simple-restful-api/cmd/config"
	"github.com/potato/simple-restful-api/infra/db"
	"github.com/potato/simple-restful-api/pkg/domain/spec"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "query ", log.LstdFlags|log.Lshortfile)
)

func CreateQuerier(cfg *config.Query) EntityStore {
	connection := db.InitMongoOrExit(cfg.EntityStoreDB)
	entityStoreDB := connection.Collection(spec.Todo{}.TableName())
	store := NewEntityStore(entityStoreDB)
	return store
}
