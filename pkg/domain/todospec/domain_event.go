package todospec

import (
	eventModel "github.com/potato/simple-restful-api/infra/model"
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
	eventModel.EventModel
	Title  string `json:"title"`
	UserNo uint   `json:"user_no"`
}

type TitleUpdated struct {
	eventModel.EventModel
	Title  string `json:"title"`
	UserNo uint   `json:"user_no"`
}

type StatusUpdated struct {
	eventModel.EventModel
	Status string `json:"status"`
	UserNo uint   `json:"user_no"`
}

type TodoDeleted struct {
	eventModel.EventModel
	AppliedAt time.Time `json:"applied_at"`
	UserNo    uint      `json:"user_no"`
}
