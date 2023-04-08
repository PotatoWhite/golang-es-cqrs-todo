package todospec

import (
	"github.com/easywalk/go-simply-cqrs"
	"github.com/google/uuid"
	"time"
)

// NewTodoCreatedEvent create TodoCreated command
func NewTodoCreatedEvent(aggregateID uuid.UUID, userNo uint, title string) TodoCreated {
	return TodoCreated{
		EventModel: simply.EventModel{
			AggregateID: aggregateID,
			EventType:   TodoCreatedEvent,
			AppliedAt:   time.Now(),
		},
		Title:  title,
		UserNo: userNo,
	}
}

// NewTitleUpdatedEvent create TitleUpdated command
func NewTitleUpdatedEvent(aggregateID uuid.UUID, userNo uint, title string) TitleUpdated {
	return TitleUpdated{
		EventModel: simply.EventModel{
			AggregateID: aggregateID,
			EventType:   TitleUpdatedEvent,
			AppliedAt:   time.Now(),
		},
		Title:  title,
		UserNo: userNo,
	}
}

// NewStatusUpdatedEvent create StatusUpdated command
func NewStatusUpdatedEvent(aggregateID uuid.UUID, userNo uint, status string) StatusUpdated {
	return StatusUpdated{
		EventModel: simply.EventModel{
			AggregateID: aggregateID,
			EventType:   StatusUpdatedEvent,
			AppliedAt:   time.Now(),
		},
		Status: status,
		UserNo: userNo,
	}
}

// NewTodoDeletedEvent create TodoDeleted command
func NewTodoDeletedEvent(aggregateID uuid.UUID) TodoDeleted {
	return TodoDeleted{
		EventModel: simply.EventModel{
			AggregateID: aggregateID,
			EventType:   TodoDeletedEvent,
			AppliedAt:   time.Now(),
		},
	}
}
