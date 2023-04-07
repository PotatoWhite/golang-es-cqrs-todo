package projector

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/potato/simple-restful-api/cmd/config"
	"github.com/potato/simple-restful-api/infra/command"
	"github.com/potato/simple-restful-api/infra/db"
	"github.com/potato/simple-restful-api/infra/model"
	"github.com/potato/simple-restful-api/infra/query"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "projector ", log.LstdFlags|log.Lshortfile)
)

func InitProjector(cfg *config.Projector, ec *chan eventModel.Event) {
	p, _ := NewObserver(cfg.Kafka, ec)

	initTokenStore("todo", cfg)

	go p.Run()
}

func initTokenStore(targetGroup string, cfg *config.Projector) TokenStore {
	tokenStoreDB, err := db.InitPostgresOrExit(cfg.TokenStoreDB)
	if err != nil {
		logger.Fatalln("Error initializing token store database", err)
	}

	// auto migration for TokenStore
	if err := tokenStoreDB.Migrator().DropTable(&command.Token{}); err != nil {
		logger.Fatalln("Error dropping table", err)
	}
	logger.Printf("Dropped table %s", command.Token{}.TableName())

	if err := tokenStoreDB.AutoMigrate(&command.Token{}); err != nil {
		logger.Fatalln("Error auto migrating database", err)
	}
	logger.Printf("Created tables of TokenStore : %s", command.Token{}.TableName())

	// token store
	tokenStore := NewTokenStore(tokenStoreDB)
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

func RunConsumer(evs command.EventStore, ets query.EntityStore, cfg *config.KafkaConfig) {
	kCfg := sarama.NewConfig()
	kCfg.Producer.Return.Successes = true
	conString := []string{cfg.BootstrapServers}
	conn, err := sarama.NewClient(conString, kCfg)
	if err != nil {
		logger.Fatalf("Error creating client: %v", err)
	}

	consumer, err := sarama.NewConsumerGroupFromClient(cfg.ConsumerGroup, conn)
	if err != nil {
		logger.Fatalf("Error creating consumer group client: %v", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Fatalf("Error closing consumer: %v", err)
		}
	}()

	ctx := context.Background()

	p := &ProjectHandler{
		Evs: evs,
		Ets: ets,
	}

	for {
		err := consumer.Consume(ctx, []string{cfg.Topic}, p)
		if err != nil {
			logger.Fatalf("Error from consumer: %v", err)
		}
	}
}
