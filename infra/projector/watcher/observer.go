package watcher

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/potato/simple-restful-api/infra/config"
	"github.com/potato/simple-restful-api/infra/model"
)

func NewObserver(cfg *config.KafkaConfig, ec <-chan eventModel.Event) (Observer, error) {
	// initialize kafka producer
	producer, err := sarama.NewSyncProducer([]string{cfg.BootstrapServers}, nil)
	if err != nil {
		logger.Fatalln("Error initializing kafka", err)
		return nil, err
	}

	logger.Println("Kafka Producer initialized")

	var p = &observer{
		producer: producer,
		ec:       ec,
		topic:    cfg.Topic,
	}
	p.Run()

	return p, nil
}

type Observer interface {
	Run()
}

type observer struct {
	ec       <-chan eventModel.Event
	producer sarama.SyncProducer
	topic    string
}

func (p *observer) Run() {
	// event loop
	go func() {
		for {
			select {
			case event := <-p.ec:
				p.publishEvent(&event)
			}
		}
	}()
}

func (p *observer) publishEvent(event *eventModel.Event) {
	// publish event to kafka
	jsonPayload, err := json.Marshal(event)
	if err != nil {
		logger.Println("error marshalling command payload: ", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("eventType"),
				Value: []byte((*event).Type()),
			},
		},
		Key:   sarama.StringEncoder((*event).ID().String()),
		Value: sarama.StringEncoder(jsonPayload),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		logger.Println("error publishing event: ", err)
	}
}
