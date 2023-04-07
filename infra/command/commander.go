package command

import (
	"github.com/potato/simple-restful-api/infra/config"
	"github.com/potato/simple-restful-api/infra/db"
	"github.com/potato/simple-restful-api/infra/model"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "commander: ", log.LstdFlags|log.Lshortfile)
)

func CreateCommander(cfg *config.Command, ec chan eventModel.Event) EventStore {
	eventStoreDB, err := db.InitPostgresOrExit(cfg.EventStoreDB)
	if err != nil {
		logger.Fatalln("Error connecting to database", err)
	}
	eventStore := createEventStoreOrExit(eventStoreDB, ec)
	return eventStore
}

func createEventStoreOrExit(eventStoreDB *gorm.DB, ec chan eventModel.Event) EventStore {
	// auto migration for EventStore
	if err := eventStoreDB.Migrator().DropTable(&Event{}); err != nil {
		logger.Fatalln("Error dropping table", err)
	}
	logger.Printf("Dropped table %s", Event{}.TableName())

	if err := eventStoreDB.AutoMigrate(&Event{}); err != nil {
		logger.Fatalln("Error auto migrating database", err)
	}
	logger.Printf("Created tables of EventStore : %s", Event{}.TableName())

	// command store
	eventStore := NewEventStore(eventStoreDB, ec)
	logger.Println("Event store initialized")

	return eventStore
}
