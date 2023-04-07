package spec

import (
	"github.com/google/uuid"
	"time"
)

// Events enum
const (
	TodoCreatedEvent   = "TodoCreatedEvent"
	TitleUpdatedEvent  = "TitleUpdatedEvent"
	StatusUpdatedEvent = "StatusUpdatedEvent"
	TodoDeletedEvent   = "TodoDeletedEvent"
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

type TodoCreated struct {
	EventModel
	Title  string `json:"title"`
	UserNo uint   `json:"user_no"`
}

type TitleUpdated struct {
	EventModel
	Title  string `json:"title"`
	UserNo uint   `json:"user_no"`
}

type StatusUpdated struct {
	EventModel
	Status string `json:"status"`
	UserNo uint   `json:"user_no"`
}

type TodoDeleted struct {
	EventModel
	AppliedAt time.Time `json:"applied_at"`
	UserNo    uint      `json:"user_no"`
}
