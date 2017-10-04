package main

import "fmt"

type DeployCommand struct {
	AbstractCommand
	Positional struct {
		Server string `description:"server configuration key"`
	} `positional-args:"true"`
	Dump  bool   `long:"dump"  description:"dump configuration as yaml"`
}

// Run deployment command
func (command *DeployCommand) Execute(args []string) error {
	config := command.GetConfig()
	server := command.getServerSelectionFromUser(config, "deploy", command.Positional.Server)
	confServer, err := config.GetDeployServer(server)
	if err != nil {
		Logger.FatalErrorExit(3, err)
	}
	Logger.Step("using Server[%s]", server)
	Logger.Step("using %s", confServer.Connection.String())

	// --dump
	if command.Dump {
		fmt.Println()
		fmt.Println(confServer.AsYaml())
	} else {
		confServer.Deploy()
		Logger.Println("-> finished")
	}

	return nil
}
