package watcher

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/potato/simple-restful-api/infra/command"
	"github.com/potato/simple-restful-api/infra/config"
	"github.com/potato/simple-restful-api/infra/model"
	"github.com/potato/simple-restful-api/infra/projector/generator"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "watcher ", log.LstdFlags|log.Lshortfile)
)

func CreateProjector(cfg *config.Projector, ec <-chan eventModel.Event) Observer {
	p, err := NewObserver(cfg.Kafka, ec)
	if err != nil {
		logger.Fatalln("Error initializing projector", err)
	}
	return p
}

func StartEntityGenerator(evs command.EventStore, eg generator.EntityGenerator, cfg *config.KafkaConfig) {
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

	p := &transformer{
		Evs: evs,
		Eg:  eg,
	}

	for {
		err := consumer.Consume(ctx, []string{cfg.Topic}, p)
		if err != nil {
			logger.Fatalf("Error from consumer: %v", err)
		}
	}
}
