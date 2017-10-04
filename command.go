package main

import (
	"strings"
	"./sync"
)

type AbstractCommand struct {

}

func (command *AbstractCommand) GetConfig() *sync.SyncConfig {
	Logger.Main("Initialisation")
	configFile := findConfigFile()
	if configFile == "" {
		Logger.FatalExit(2, "Unable to find configuration file (searched  %s)", strings.Join(validConfigFiles, " "))
	}
	Logger.Step("found configuration file %s", configFile)

	sync.Logger = Logger
	config := sync.NewConfigParser(configFile)

	return config
}
