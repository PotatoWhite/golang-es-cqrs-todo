package spec

import (
	"github.com/google/uuid"
	eventModel "github.com/potato/simple-restful-api/infra/model"
	"time"
)

// create TodoCreated command
func NewTodoCreatedEvent(aggregateID uuid.UUID, userNo uint, title string) TodoCreated {

	return TodoCreated{
		EventModel: eventModel.EventModel{
			AggregateID: aggregateID,
			EventType:   TodoCreatedEvent,
			AppliedAt:   time.Now(),
		},
		Title:  title,
		UserNo: userNo,
	}
}

// create TitleUpdated command
func NewTitleUpdatedEvent(aggregateID uuid.UUID, userNo uint, title string) TitleUpdated {
	return TitleUpdated{
		EventModel: eventModel.EventModel{
			AggregateID: aggregateID,
			EventType:   TitleUpdatedEvent,
			AppliedAt:   time.Now(),
		},
		Title:  title,
		UserNo: userNo,
	}
}

// create StatusUpdated command
func NewStatusUpdatedEvent(aggregateID uuid.UUID, userNo uint, status string) StatusUpdated {
	return StatusUpdated{
		EventModel: eventModel.EventModel{
			AggregateID: aggregateID,
			EventType:   StatusUpdatedEvent,
			AppliedAt:   time.Now(),
		},
		Status: status,
		UserNo: userNo,
	}
}

// create TodoDeleted command
func NewTodoDeletedEvent(aggregateID uuid.UUID) TodoDeleted {
	return TodoDeleted{
		EventModel: eventModel.EventModel{
			AggregateID: aggregateID,
			EventType:   TodoDeletedEvent,
			AppliedAt:   time.Now(),
		},
	}
}
