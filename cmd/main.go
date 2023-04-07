package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/potato/simple-restful-api/cmd/config"
	"github.com/potato/simple-restful-api/infra/command"
	"github.com/potato/simple-restful-api/infra/model"
	"github.com/potato/simple-restful-api/infra/projector"
	"github.com/potato/simple-restful-api/infra/query"
	"github.com/potato/simple-restful-api/pkg/rest/command"
	"github.com/potato/simple-restful-api/pkg/rest/query"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "main ", log.LstdFlags|log.Lshortfile)
)

func main() {
	cfg := loadCfg()

	eventChannel := make(chan eventModel.Event, 100)

	eventStore := command.CreateCommander(cfg.Command, &eventChannel)
	entityStore := query.CreateQuerier(cfg.Query)

	projector.InitProjector(cfg.Projector, &eventChannel)
	go projector.RunConsumer(eventStore, entityStore, cfg.Projector.Kafka)

	engine := gin.Default()
	todoGroup := engine.Group("")
	commandApi.NewTodoRouter(todoGroup, eventStore)
	queryApi.NewTodoRouter(todoGroup, entityStore)

	if err := engine.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		logger.Fatalln("Error running server", err)
	}
}

func loadCfg() config.Config {
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
