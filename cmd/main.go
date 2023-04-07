package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/potato/simple-restful-api/cmd/config"
	"github.com/potato/simple-restful-api/pkg/domain/command"
	"github.com/potato/simple-restful-api/pkg/domain/projector"
	"github.com/potato/simple-restful-api/pkg/domain/query"
	"github.com/potato/simple-restful-api/pkg/domain/spec"
	"github.com/potato/simple-restful-api/pkg/rest/command"
	"github.com/potato/simple-restful-api/pkg/rest/query"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var (
	logger = log.New(os.Stdout, "main ", log.LstdFlags|log.Lshortfile)
)

func main() {
	cfg := LoadCfg()

	eventChannel := make(chan spec.Event, 100)
	// command
	eventStore := createCommander(cfg.Command, &eventChannel)
	entityStore := createQuerier(cfg.Query)

	// projector
	initProjector(cfg.Projector, &eventChannel)

	go runConsumer(eventStore, entityStore, cfg.Projector.Kafka)

	// query
	//readModelStore := createQuerier(cfg, eventStore)

	// initialize gin
	engine := gin.Default()
	todoGroup := engine.Group("")
	registerCommandRoutes(todoGroup, eventStore)
	registerQueryRoutes(todoGroup, entityStore)

	// run gin
	if err := engine.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		logger.Fatalln("Error running server", err)
	}

	//initCommandHandlerOrExit(eventStore, cfg.CmdDB.AutoMigration)
}

func createCommander(cfg *config.Command, ec *chan spec.Event) command.EventStore {
	eventStoreDB := initCommandDatabasesOrExit(cfg.EventStoreDB)
	eventStore := createEventStoreOrExit(eventStoreDB, ec)
	return eventStore
}

func createQuerier(cfg *config.Query) query.EntityStore {
	entityStoreDB := initMongoOrExit(cfg)
	store := query.NewEntityStore(entityStoreDB)
	return store
}

func initMongoOrExit(cfg *config.Query) (entityStore *mongo.Collection) {
	// create mongo connection
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.EntityStoreDB.Host, cfg.EntityStoreDB.Port)).SetAuth(
		options.Credential{
			Username:      cfg.EntityStoreDB.User,
			Password:      cfg.EntityStoreDB.Password,
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    cfg.EntityStoreDB.Database,
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

	return client.Database(cfg.EntityStoreDB.Database).Collection(spec.Todo{}.TableName())
}

func initProjector(cfg *config.Projector, ec *chan spec.Event) {
	// init Producer
	p, _ := projector.NewObserver(cfg.Kafka, ec)

	// init Consumer
	initTokenStore("todo", cfg)

	// create projector

	go p.Run()

}

func runConsumer(evs command.EventStore, ets query.EntityStore, cfg *config.KafkaConfig) {
	kcfg := sarama.NewConfig()
	kcfg.Producer.Return.Successes = true
	conString := []string{cfg.BootstrapServers}
	conn, err := sarama.NewClient(conString, kcfg)
	if err != nil {
		logger.Fatalf("Error creating client: %v", err)
	}

	// Create a new consumer group
	consumer, err := sarama.NewConsumerGroupFromClient(spec.Todo{}.TableName(), conn)
	if err != nil {
		logger.Fatalf("Error creating consumer group client: %v", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Fatalf("Error closing consumer: %v", err)
		}
	}()

	ctx := context.Background()

	// Iterate over consumer sessions.
	for {
		err := consumer.Consume(ctx, []string{cfg.Topic}, &projector.ProjectHandler{
			Evs: evs,
			Ets: ets,
		})
		if err != nil {
			logger.Fatalf("Error from consumer: %v", err)
		}
	}
}

func initTokenStore(targetGroup string, cfg *config.Projector) command.TokenStore {
	eventStoreDB := initCommandDatabasesOrExit(cfg.TokenStoreDB)

	// auto migration for TokenStore
	if err := eventStoreDB.Migrator().DropTable(&command.Token{}); err != nil {
		logger.Fatalln("Error dropping table", err)
	}
	logger.Printf("Dropped table %s", command.Token{}.TableName())

	if err := eventStoreDB.AutoMigrate(&command.Token{}); err != nil {
		logger.Fatalln("Error auto migrating database", err)
	}
	logger.Printf("Created tables of TokenStore : %s", command.Token{}.TableName())

	// token store
	tokenStore := command.NewTokenStore(eventStoreDB)
	logger.Println("Token store initialized")

	// logging last token
	if lastToken, err := tokenStore.GetLastTokenByTargetGroup(targetGroup); err != nil {
		logger.Println("Error getting last token", err)
		// make new token
		if token, err := tokenStore.CreateToken(targetGroup, 0); err != nil {
			logger.Println("Error creating token", err)
		} else {
			logger.Println("Created token", token)
		}
	} else {
		logger.Println("Last token :", lastToken)
	}

	return tokenStore
}

func LoadCfg() config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalln("Error loading config", err)
		return config.Config{}
	}

	// print all config
	_config, err := json.Marshal(cfg)
	logger.Println("Config loaded successfully : ", string(_config))
	return cfg
}

func initCommandDatabasesOrExit(cfg *config.DBConfig) (eventStoreDB *gorm.DB) {
	// initialize database of EventStore
	eventDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul", cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port)
	eventStoreDB, err := gorm.Open(postgres.Open(eventDsn), &gorm.Config{})

	if err != nil {
		logger.Fatalln("Error connecting to database", err)
		return nil
	}

	logger.Println("Database initialized")

	return eventStoreDB
}

func createEventStoreOrExit(eventStoreDB *gorm.DB, ec *chan spec.Event) command.EventStore {
	// auto migration for EventStore
	if err := eventStoreDB.Migrator().DropTable(&command.Event{}); err != nil {
		logger.Fatalln("Error dropping table", err)
	}
	logger.Printf("Dropped table %s", command.Event{}.TableName())

	if err := eventStoreDB.AutoMigrate(&command.Event{}); err != nil {
		logger.Fatalln("Error auto migrating database", err)
	}
	logger.Printf("Created tables of EventStore : %s", command.Event{}.TableName())

	// command store
	eventStore := command.NewEventStore(eventStoreDB, ec)
	logger.Println("Event store initialized")

	return eventStore
}

func registerCommandRoutes(group *gin.RouterGroup, evs command.EventStore) {
	commandApi.NewTodoRouter(group, evs)
}

func registerQueryRoutes(group *gin.RouterGroup, ets query.EntityStore) {
	queryApi.NewTodoRouter(group, ets)
}
