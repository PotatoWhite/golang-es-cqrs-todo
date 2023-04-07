package projector

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/potato/simple-restful-api/infra/command"
	"github.com/potato/simple-restful-api/infra/model"
	"github.com/potato/simple-restful-api/infra/query"
	"github.com/potato/simple-restful-api/pkg/domain/spec"
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
		var event eventModel.EventModel
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
	todo, err := p.Ets.ReplayEvents(events)
	if err != nil {
		return todo, err
	}

	return todo, nil
}
