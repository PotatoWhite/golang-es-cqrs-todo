package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/potato/simple-restful-api/infra/command"
	"github.com/potato/simple-restful-api/infra/config"
	"github.com/potato/simple-restful-api/infra/model"
	"github.com/potato/simple-restful-api/infra/monitoring"
	"github.com/potato/simple-restful-api/infra/projector/watcher"
	"github.com/potato/simple-restful-api/pkg"
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

	// init event channel
	eventChannel := make(chan eventModel.Event, 100)

	// start event store
	eventStore := command.CreateCommander(cfg.Command, eventChannel) // events -> event store
	todoStore := todo.CreateTodoService(cfg.Query)                   // events -> entity -> pkg store

	generator := todo.NewEntityGenerator(todoStore)                             // events -> entity
	go watcher.CreateProjector(cfg.Projector, eventChannel).Run()               // event channel -> kafka
	go watcher.StartEntityGenerator(eventStore, generator, cfg.Projector.Kafka) // kafka -> events -> projector

	engine := gin.Default()
	todoGroup := engine.Group("")

	commandApi.NewTodoRouter(todoGroup, eventStore)
	queryApi.NewTodoRouter(todoGroup, todoStore)

	engine.Use(monitoring.LoggingMiddleware())

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
