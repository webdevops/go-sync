package main

import (
	"strings"
	"os"
	"fmt"
	"path"
	"./sync"
	"gopkg.in/AlecAivazis/survey.v1"

)

type AbstractCommand struct {

}

func (command *AbstractCommand) GetConfig() *sync.SyncConfig {
	Logger.Main("Initialisation")
	configFile := command.findConfigFile()
	if configFile == "" {
		Logger.FatalExit(2, "Unable to find configuration file (searched  %s)", strings.Join(validConfigFiles, " "))
	}
	Logger.Step("found configuration file %s", configFile)

	sync.Logger = Logger
	config := sync.NewConfigParser(configFile)

	return config
}

func (command *AbstractCommand) findConfigFile() string {
	pwd, err := os.Getwd()
	if err != nil {
		Logger.FatalErrorExit(1, err)
		fmt.Println(err)
	}

	for true {
		for _, filename := range validConfigFiles {
			filepath := path.Join(pwd, filename)
			if sync.FileExists(filepath) {
				return filepath
			}
		}


		// already found root, we finished here
		if pwd == "/" {
			break
		}

		pwd = path.Dir(pwd)

		// oh, path seems to be empty.. stopping here now
		if pwd == "." || pwd == "" {
			break
		}
	}

	return ""
}


func (command *AbstractCommand) getServerSelectionFromUser(config *sync.SyncConfig, confType string, userSelection string) string {
	if userSelection == "" {
		prompt := &survey.Select{
			Message: "Choose configuration:",
			Options: config.GetServerList(confType),
		}
		survey.AskOne(prompt, &userSelection, nil)
	}

	return userSelection
}
