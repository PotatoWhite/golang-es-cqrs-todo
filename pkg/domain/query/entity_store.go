package query

import (
	"context"
	"github.com/google/uuid"
	"github.com/potato/simple-restful-api/pkg/domain/spec"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "query ", log.LstdFlags|log.Lshortfile)
)

func NewEntityStore(entityCollection *mongo.Collection) EntityStore {
	return &entityStore{
		collection: entityCollection,
	}
}

type EntityStore interface {
	GetTodo(id uuid.UUID) (*spec.Todo, error)
	GetTodosByUserNo(userNo uint) ([]spec.Todo, error)
	GetTodosByUserNoAndStatus(userNo uint, status string) ([]spec.Todo, error)
	GetTodosByUserNoAndNotStatus(userNo uint, status string) ([]spec.Todo, error)
	SaveTodo(todo *spec.Todo) error
}

type entityStore struct {
	collection *mongo.Collection
}

func (ets *entityStore) SaveTodo(todo *spec.Todo) error {
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

	logger.Println("Saved todo", todo.Id)
	return nil
}

func (ets *entityStore) GetTodo(id uuid.UUID) (todo *spec.Todo, err error) {
	filter := bson.M{"id": id}
	return ets.getTodoByFilter(filter)
}

func (ets *entityStore) GetTodosByUserNo(userNo uint) ([]spec.Todo, error) {
	filter := bson.M{"user_no": userNo}
	return ets.getAllTodosByFilter(filter)
}

func (ets *entityStore) GetTodosByUserNoAndStatus(userNo uint, status string) ([]spec.Todo, error) {
	filter := bson.M{"user_no": userNo, "status": status}
	return ets.getAllTodosByFilter(filter)
}

func (ets *entityStore) GetTodosByUserNoAndNotStatus(userNo uint, status string) ([]spec.Todo, error) {
	filter := bson.M{"user_no": userNo, "status": bson.M{"$ne": status}}
	return ets.getAllTodosByFilter(filter)
}

func (ets *entityStore) getAllTodosByFilter(filter bson.M) ([]spec.Todo, error) {
	var todos []spec.Todo
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
		var todo spec.Todo
		if err := cs.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (ets *entityStore) getTodoByFilter(filter bson.M) (*spec.Todo, error) {
	var todo spec.Todo
	if err := ets.collection.FindOne(context.Background(), filter).Decode(&todo); err != nil {
		return nil, err
	}
	return &todo, nil
}

func toBSON(todo *spec.Todo) bson.M {
	return bson.M{
		"id":         todo.Id,
		"user_no":    todo.UserNo,
		"title":      todo.Title,
		"status":     todo.Status,
		"applied_at": todo.AppliedAt,
	}
}
