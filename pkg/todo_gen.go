package todo

import (
	"github.com/potato/simple-restful-api/infra/command"
	"github.com/potato/simple-restful-api/infra/config"
	"github.com/potato/simple-restful-api/infra/db"
	"github.com/potato/simple-restful-api/infra/projector/generator"
	"github.com/potato/simple-restful-api/pkg/domain/todospec"
	"github.com/potato/simple-restful-api/pkg/repository"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "query ", log.LstdFlags|log.Lshortfile)
)

func CreateTodoService(cfg *config.Query) repository.TodoStore {
	connection := db.InitMongoOrExit(cfg.EntityStoreDB)
	entityStoreDB := connection.Collection(todospec.Todo{}.TableName())
	store := repository.NewTodoStore(entityStoreDB)
	return store
}

func NewEntityGenerator(ets repository.TodoStore) generator.EntityGenerator {
	eg := &todoGenerator{
		Ets: ets,
	}
	return eg
}

type todoGenerator struct {
	Ets repository.TodoStore
}

func (eg *todoGenerator) CreateEntityAnsSave(events []*command.Event) error {
	// calculate pkg entity and update query store
	todo, err := eg.calculateTodoEntity(events)
	if err != nil {
		return err
	}

	if err := eg.Ets.SaveTodo(todo); err != nil {
		logger.Println("Error saving pkg", err)
		return err
	}
	return nil
}

func (eg *todoGenerator) calculateTodoEntity(events []*command.Event) (*todospec.Todo, error) {
	logger.Println("events size : ", len(events))
	todo, err := eg.Ets.ReplayEvents(events)
	if err != nil {
		return todo, err
	}

	return todo, nil
}
