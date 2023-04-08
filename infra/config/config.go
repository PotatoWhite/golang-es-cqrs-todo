package config

import (
	"github.com/easywalk/go-simply-cqrs"
	"github.com/spf13/viper"
	"log"
	"os"
)

// for viper
var (
	logger = log.New(os.Stdout, "config ", log.LstdFlags|log.Lshortfile)
)

type Config struct {
	App       *AppConfig
	Command   *simply.Command
	Query     *simply.Query
	Projector *simply.Projector
	Server    *ServerConfig
	Profile   *profile
}

type AppConfig struct {
	Name string
}

type ServerConfig struct {
	Port string
}

type profile struct {
	Name string
}

func LoadConfig() (cfg Config, e error) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./cmd")
	viper.AddConfigPath("./cmd/config")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// from kubernetes secret
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatalln("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Fatalln("Unable to decode into struct, %v", err)
	}

	logger.Println("Config loaded")

	return cfg, nil
}
