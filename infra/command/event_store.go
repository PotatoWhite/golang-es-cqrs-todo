package command

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/potato/simple-restful-api/infra/model"
	"gorm.io/gorm"
)

func NewEventStore(db *gorm.DB, ec chan eventModel.Event) EventStore {
	return &eventStore{
		db:           db,
		eventChannel: ec,
	}
}

type EventStore interface {
	AddAndPublishEvent(userNo uint, event eventModel.Event) (*Event, error)
	GetAllEvents(aggregateId uuid.UUID) ([]*Event, error)
	GetLastEvent(aggregateId uuid.UUID) (*Event, error)
}

type eventStore struct {
	db           *gorm.DB
	eventChannel chan eventModel.Event
}

func (evs *eventStore) GetLastEvent(aggregateId uuid.UUID) (event *Event, err error) {
	return event, evs.db.Where("aggregate_id = ?", aggregateId).Last(&event).Error
}

func (evs *eventStore) GetAllEvents(aggregateId uuid.UUID) (events []*Event, err error) {
	return events, evs.db.Where("aggregate_id = ?", aggregateId).Find(&events).Error
}

func (evs *eventStore) AddAndPublishEvent(userNo uint, event eventModel.Event) (eventEntity *Event, err error) {

	jsonPayload, err := json.Marshal(event)
	if err != nil {
		logger.Println("error marshalling command payload: ", err)
	}

	eventEntity = &Event{
		UserNo:      userNo,
		EventType:   event.Type(),
		AggregateId: event.ID(),
		Payload:     jsonPayload,
	}

	err = evs.db.Create(eventEntity).Error
	if err != nil {
		return nil, err
	}

	// publish event to Projector by channel if it is set
	if evs.eventChannel != nil {
		evs.eventChannel <- event
	}

	return eventEntity, nil
}
