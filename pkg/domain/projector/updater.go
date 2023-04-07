package projector

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/potato/simple-restful-api/pkg/domain/command"
	"github.com/potato/simple-restful-api/pkg/domain/query"
	"github.com/potato/simple-restful-api/pkg/domain/spec"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "observer ", log.LstdFlags|log.Lshortfile)
)

type ProjectHandler struct {
	Evs command.EventStore
	Ets query.EntityStore
	//ReadStore query.ReadStore
}

func (p *ProjectHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (p *ProjectHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (p *ProjectHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// unmarshal msg.Value by eventType(Header:eventType)
		var event spec.EventModel
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			// send to dead letter queue and continue

		} else {
			// calculate todo entity and update query store
			todo, err := p.calculateTodoEntity(event.ID())
			//saved := p.updateTodoReadStore(todo)

			if err := p.Ets.SaveTodo(todo); err != nil {
				logger.Println("Error saving todo", err)
				sess.MarkMessage(msg, "DLQ")
				return err
			}
			if err != nil {
				logger.Println("Error creating todo", err)
			}
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (p *ProjectHandler) calculateTodoEntity(aggregateId uuid.UUID) (*spec.Todo, error) {
	logger.Println("calculateTodoEntity : ", aggregateId)
	// get all events
	events, err := p.Evs.GetAllEvents(aggregateId)
	if err != nil {
		return nil, err
	}
	logger.Println("events size : ", len(events))

	todo, err := p.replay(events)
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (p *ProjectHandler) replay(events []*command.Event) (todo *spec.Todo, err error) {
	for _, event := range events {
		switch event.EventType {
		case spec.TodoCreatedEvent:
			var todoCreated spec.TodoCreated
			err := json.Unmarshal(event.Payload, &todoCreated)
			todo, err = p.createTodo(todoCreated)
			if err != nil {
				return nil, err
			}

		case spec.TitleUpdatedEvent:
			var titleUpdated spec.TitleUpdated
			err := json.Unmarshal(event.Payload, &titleUpdated)
			todo.Title = titleUpdated.Title
			todo.AppliedAt = titleUpdated.AppliedAt
			if err != nil {
				return nil, err
			}

		case spec.StatusUpdatedEvent:
			var statusUpdated spec.StatusUpdated
			err := json.Unmarshal(event.Payload, &statusUpdated)
			todo.Status = statusUpdated.Status
			todo.AppliedAt = statusUpdated.AppliedAt
			if err != nil {
				return nil, err
			}

		case spec.TodoDeletedEvent:
			var todoDeleted spec.TodoDeleted
			err := json.Unmarshal(event.Payload, &todoDeleted)
			todo.Status = spec.TodoStatusDeleted
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

func (p *ProjectHandler) createTodo(event spec.TodoCreated) (todo *spec.Todo, err error) {
	return &spec.Todo{
		Id:        event.ID(),
		Title:     event.Title,
		Status:    spec.TodoStatusInProgress,
		AppliedAt: event.AppliedAt,
	}, nil
}
