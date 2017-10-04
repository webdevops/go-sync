package main

import "fmt"

type SyncCommand struct {
	AbstractCommand
	Positional struct {
		Server  string `description:"server configuration key"`
	} `positional-args:"true"`
	Dump  bool   `long:"dump"  description:"dump configuration as yaml"`
}

// Run sync command
func (command *SyncCommand) Execute(args []string) error {
	config := command.GetConfig()
	server := command.getServerSelectionFromUser(config, "sync", command.Positional.Server)
	confServer, err := config.GetSyncServer(server)
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
		confServer.Sync()
		Logger.Println("-> finished")
	}

	return nil
}
