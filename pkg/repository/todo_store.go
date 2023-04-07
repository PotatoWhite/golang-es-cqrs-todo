package repository

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/potato/simple-restful-api/infra/command"
	"github.com/potato/simple-restful-api/infra/projector/generator"
	"github.com/potato/simple-restful-api/pkg/domain/todospec"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "repository ", log.LstdFlags|log.Lshortfile)
)

func NewTodoStore(entityCollection *mongo.Collection) TodoStore {
	return &todoStore{
		collection: entityCollection,
	}
}

type TodoStore interface {
	GetTodo(id uuid.UUID) (*todospec.Todo, error)
	GetTodoByUserNoAndId(userNo uint, id uuid.UUID) (*todospec.Todo, error)
	GetTodosByUserNo(userNo uint) ([]todospec.Todo, error)
	GetTodosByUserNoAndStatus(userNo uint, status string) ([]todospec.Todo, error)
	GetTodosByUserNoAndNotStatus(userNo uint, status string) ([]todospec.Todo, error)
	SaveTodo(todo *todospec.Todo) error
	ReplayEvents(events []*command.Event) (todo *todospec.Todo, err error)
}

type todoStore struct {
	collection *mongo.Collection
	Eg         generator.EntityGenerator
}

func (ets *todoStore) GetTodoByUserNoAndId(userNo uint, id uuid.UUID) (*todospec.Todo, error) {
	filter := bson.M{"user_no": userNo, "id": id}
	return ets.getTodoByFilter(filter)
}

func (ets *todoStore) SaveTodo(todo *todospec.Todo) error {
	if todo == nil {
		return nil
	}

	filter := bson.M{"id": todo.Id}
	update := bson.M{"$set": toBSON(todo)}

	opts := options.Update().SetUpsert(true)

	_, err := ets.collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}

	logger.Println("Saved pkg", todo.Id)
	return nil
}

func (ets *todoStore) GetTodo(id uuid.UUID) (todo *todospec.Todo, err error) {
	filter := bson.M{"id": id}
	return ets.getTodoByFilter(filter)
}

func (ets *todoStore) GetTodosByUserNo(userNo uint) ([]todospec.Todo, error) {
	filter := bson.M{"user_no": userNo}
	return ets.getAllTodosByFilter(filter)
}

func (ets *todoStore) GetTodosByUserNoAndStatus(userNo uint, status string) ([]todospec.Todo, error) {
	filter := bson.M{"user_no": userNo, "status": status}
	return ets.getAllTodosByFilter(filter)
}

func (ets *todoStore) GetTodosByUserNoAndNotStatus(userNo uint, status string) ([]todospec.Todo, error) {
	filter := bson.M{"user_no": userNo, "status": bson.M{"$ne": status}}
	return ets.getAllTodosByFilter(filter)
}

func (ets *todoStore) getAllTodosByFilter(filter bson.M) ([]todospec.Todo, error) {
	var todos []todospec.Todo
	cs, err := ets.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cs *mongo.Cursor, ctx context.Context) {
		err := cs.Close(ctx)
		if err != nil {
			logger.Println("Error closing cursor", err)
		}
	}(cs, context.Background())
	for cs.Next(context.Background()) {
		var todo todospec.Todo
		if err := cs.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (ets *todoStore) getTodoByFilter(filter bson.M) (*todospec.Todo, error) {
	var todo todospec.Todo
	if err := ets.collection.FindOne(context.Background(), filter).Decode(&todo); err != nil {
		return nil, err
	}
	return &todo, nil
}

func toBSON(todo *todospec.Todo) bson.M {
	return bson.M{
		"id":         todo.Id,
		"user_no":    todo.UserNo,
		"title":      todo.Title,
		"status":     todo.Status,
		"applied_at": todo.AppliedAt,
	}
}

func (ets *todoStore) ReplayEvents(events []*command.Event) (todo *todospec.Todo, err error) {
	for _, event := range events {
		switch event.EventType {
		case todospec.TodoCreatedEvent:
			var todoCreated todospec.TodoCreated
			err := json.Unmarshal(event.Payload, &todoCreated)
			todo, err = ets.createTodo(todoCreated)
			if err != nil {
				return nil, err
			}

		case todospec.TitleUpdatedEvent:
			var titleUpdated todospec.TitleUpdated
			err := json.Unmarshal(event.Payload, &titleUpdated)
			todo.Title = titleUpdated.Title
			todo.AppliedAt = titleUpdated.AppliedAt
			if err != nil {
				return nil, err
			}

		case todospec.StatusUpdatedEvent:
			var statusUpdated todospec.StatusUpdated
			err := json.Unmarshal(event.Payload, &statusUpdated)
			todo.Status = statusUpdated.Status
			todo.AppliedAt = statusUpdated.AppliedAt
			if err != nil {
				return nil, err
			}

		case todospec.TodoDeletedEvent:
			var todoDeleted todospec.TodoDeleted
			err := json.Unmarshal(event.Payload, &todoDeleted)
			todo.Status = todospec.TodoStatusDeleted
			todo.AppliedAt = todoDeleted.AppliedAt
			if err != nil {
				return nil, err
			}

		default:
			logger.Printf("Unknown event type: %s", event.EventType)
		}
	}
	return todo, nil
}

func (ets *todoStore) createTodo(event todospec.TodoCreated) (todo *todospec.Todo, err error) {
	return &todospec.Todo{
		Id:        event.ID(),
		Title:     event.Title,
		UserNo:    event.UserNo,
		Status:    todospec.TodoStatusInProgress,
		AppliedAt: event.AppliedAt,
	}, nil
}
