package watcher

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/potato/simple-restful-api/infra/command"
	"github.com/potato/simple-restful-api/infra/model"
	"github.com/potato/simple-restful-api/infra/projector/generator"
)

type transformer struct {
	Evs command.EventStore
	Eg  generator.EntityGenerator
}

func (p *transformer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (p *transformer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (p *transformer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// unmarshal msg.Value by eventType(Header:eventType)
		var event eventModel.EventModel
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			// send to dead letter queue and continue

		} else {
			// get all events from event store
			events, err := p.Evs.GetAllEvents(event.ID())
			if err := p.Eg.CreateEntityAnsSave(events); err != nil {
				sess.MarkMessage(msg, "Move to DLQ")
				return err
			}
			if err != nil {
				logger.Println("Error creating entity", err)
			}
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
