package command

import (
	"encoding/json"
	"github.com/google/uuid"
)

type Event struct {
	ID          uint            `gorm:"primarykey"`
	UserNo      uint            `gorm:"index"`
	EventType   string          `gorm:"index"`
	AggregateId uuid.UUID       `gorm:"type:uuid;index"`
	Payload     json.RawMessage `gorm:"type:json"`
}

func (Event) TableName() string {
	return "event_store"
}
