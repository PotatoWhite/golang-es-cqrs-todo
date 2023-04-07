package eventModel

import (
	"github.com/google/uuid"
	"time"
)

type Event interface {
	ID() uuid.UUID
	Type() string
	Time() time.Time
}

type EventModel struct {
	AggregateID uuid.UUID `json:"id" gorm:"type:uuid;column:id;index"`
	EventType   string    `json:"event_type" gorm:"index"`
	AppliedAt   time.Time `json:"applied_at"`
}

func (ev *EventModel) ID() uuid.UUID {
	return ev.AggregateID
}

func (ev *EventModel) Type() string {
	return ev.EventType
}

func (ev *EventModel) Time() time.Time {
	return ev.AppliedAt
}
