package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	ApiHost      string `envconfig:"API_HOST" default:"127.0.0.1:8000"`
	LivecheckURI string `envconfig:"LIVECHECK" default:"/live"`
}

func GetConfig() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("envconfig error: %v", err)
		return config, err
	}
	return config, nil
}
