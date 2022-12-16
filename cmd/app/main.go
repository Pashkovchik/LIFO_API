package main

import (
	"lifo-rest-api/internal/app"
	"lifo-rest-api/internal/config"
	"lifo-rest-api/pkg/configurator"
)

const (
	configFile     = "config"
	configFilePath = "configs"
)

func main() {
	conf := new(config.Config)
	configurator.InitConfigs(configFile, []string{configFilePath}, conf)

	app.Run(conf)
}
