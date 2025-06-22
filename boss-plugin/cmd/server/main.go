package main

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/boss-net/api/boss-plugin/internal/server"
	"github.com/boss-net/api/boss-plugin/internal/types/app"
	"github.com/boss-net/api/boss-plugin/internal/utils/log"
)

func main() {
	var config app.Config

	// load env
	godotenv.Load()

	err := envconfig.Process("", &config)
	if err != nil {
		log.Panic("Error processing environment variables: %s", err.Error())
	}

	config.SetDefault()

	if err := config.Validate(); err != nil {
		log.Panic("Invalid configuration: %s", err.Error())
	}

	(&server.App{}).Run(&config)
}
