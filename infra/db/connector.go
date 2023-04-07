package db

import (
	"context"
	"fmt"
	"github.com/potato/simple-restful-api/infra/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var (
	logger = log.New(os.Stdout, "db: ", log.LstdFlags|log.Lshortfile)
)

func InitPostgresOrExit(cfg *config.DBConfig) (db *gorm.DB, err error) {
	// initialize database of EventStore
	eventDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul", cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port)
	db, err = gorm.Open(postgres.Open(eventDsn), &gorm.Config{})

	if err != nil {
		logger.Fatalln("Error connecting to database", err)
		return nil, err
	}

	logger.Println("Database initialized")

	return db, nil
}

func InitMongoOrExit(cfg *config.DBConfig) (entityStore *mongo.Database) {
	// create mongo connection
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)).SetAuth(
		options.Credential{
			Username:      cfg.User,
			Password:      cfg.Password,
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    cfg.Database,
		})
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		logger.Fatalln("Error creating mongo client", err)
	}

	// connect to mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatalln("Error connecting to mongo", err)
	}

	return client.Database(cfg.Database)
}
