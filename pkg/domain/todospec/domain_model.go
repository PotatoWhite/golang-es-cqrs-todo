package todospec

import (
	"github.com/google/uuid"
	"time"
)

// enum for status
const (
	TodoStatusInProgress = "in_progress"
	TodoStatusCompleted  = "completed"
	TodoStatusCanceled   = "canceled"
	TodoStatusDeleted    = "deleted"
)

type Todo struct {
	Id        uuid.UUID `bson:"id" json:"id"`
	UserNo    uint      `bson:"user_no" json:"user_no"`
	Title     string    `bson:"title" json:"title"`
	Status    string    `bson:"status" json:"status"`
	AppliedAt time.Time `bson:"applied_at" json:"applied_at"`
}

func (Todo) TableName() string {
	return "todos"
}
