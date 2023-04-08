package todospec

import (
	"github.com/easywalk/go-simply-cqrs"
	"time"
)

// Events enum
const (
	TodoCreatedEvent   = "TodoCreatedEvent"
	TitleUpdatedEvent  = "TitleUpdatedEvent"
	StatusUpdatedEvent = "StatusUpdatedEvent"
	TodoDeletedEvent   = "TodoDeletedEvent"
)

type TodoCreated struct {
	simply.EventModel
	Title  string `json:"title"`
	UserNo uint   `json:"user_no"`
}

type TitleUpdated struct {
	simply.EventModel
	Title  string `json:"title"`
	UserNo uint   `json:"user_no"`
}

type StatusUpdated struct {
	simply.EventModel
	Status string `json:"status"`
	UserNo uint   `json:"user_no"`
}

type TodoDeleted struct {
	simply.EventModel
	AppliedAt time.Time `json:"applied_at"`
	UserNo    uint      `json:"user_no"`
}
