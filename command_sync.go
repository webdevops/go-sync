package main

import (
	"fmt"
	"github.com/webdevops/go-sync/sync"
)

type SyncCommand struct {
	AbstractCommand
	Positional struct {
		Server  string `description:"server configuration key"`
	} `positional-args:"true"`
	Dump            bool   `long:"dump"        description:"dump configuration as yaml"`
	OnlyFilesystem  bool   `long:"filesystem"  description:"sync only filesystem"`
	OnlyDatabase    bool   `long:"database"    description:"sync only database"`
	SkipExec        bool   `long:"skip-exec"   description:"skip execution"`
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
	Logger.Step("using %s", confServer.Connection.GetInstance().String())

	confServer.SetRunConfiguration(command.buildSyncRunConfig())

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

func (command *SyncCommand) buildSyncRunConfig() (conf sync.RunConfiguration) {
	// Init
	conf.Exec = true
	conf.Database = true
	conf.Filesystem = true

	if command.OnlyFilesystem {
		conf.Database = false
		conf.Filesystem = true
	}

	if command.OnlyDatabase {
		conf.Database = true
		conf.Filesystem = false
	}

	if command.SkipExec {
		conf.Exec = false
	}

	return
}
