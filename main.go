package main

import (
	"os"
	"log"
	"fmt"
	"path"
	"strings"
	"runtime/debug"
	flags "github.com/jessevdk/go-flags"
	"github.com/webdevops/go-shell"
	"gopkg.in/AlecAivazis/survey.v1"
	"./sync"
	"./logger"
)

const (
	Name    = "gosync"
	Author  = "webdevops.io"
	Version = "0.1.0"
)

var (
	Logger *logger.SyncLogger
	argparser *flags.Parser
	args []string
)

var opts struct {
	Positional struct {
		Command string `description:"What to do [help, sync, deploy or show]" choice:"show" choice:"sync" choice:"deploy" choice:"help" required:"1"`
		Server  string `description:"server configuration key"`
	} `positional-args:"true"`

	Dump               bool     `           long:"dump"                          description:"dump configuration as yaml"`
	Verbose            []bool   `short:"v"  long:"verbose"                       description:"verbose mode"`
	DryRun             bool     `           long:"dry-run"                       description:"dry run mode"`
	ShowVersion        bool     `short:"V"  long:"version"                       description:"show version and exit"`
	ShowOnlyVersion    bool     `           long:"dumpversion"                   description:"show only version number and exit"`
	ShowHelp           bool     `short:"h"  long:"help"                          description:"show this help message"`
}

var validConfigFiles = []string{
	"gosync.yml",
	"gosync.yaml",
	".gosync.yml",
	".gosync.yaml",
}

func createArgparser() {
	var err error
	argparser = flags.NewParser(&opts, flags.Default)
	args, err = argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

	// --dumpversion
	if opts.ShowOnlyVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	// --version
	if opts.ShowVersion {
		fmt.Println(fmt.Sprintf("%s version %s", Name, Version))
		fmt.Println(fmt.Sprintf("Copyright (C) 2017 %s", Author))
		os.Exit(0)
	}
}

func findConfigFile() string {
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

func getArgServer(config *sync.SyncConfig, confType string) string {
	server := opts.Positional.Server
	if server == "" {
		prompt := &survey.Select{
			Message: "Choose configuration:",
			Options: config.GetServerList(confType),
		}
		survey.AskOne(prompt, &server, nil)
	}

	return server
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println()
			if len(opts.Verbose) >= 2 {
				fmt.Println(r)
				debug.PrintStack()
			} else {
				fmt.Println(r)
			}
			os.Exit(255)
		}
	}()

	createArgparser()

	argCommand := strings.ToLower(opts.Positional.Command)

	switch {
	case len(opts.Verbose) >= 2:
		shell.Trace = true
		shell.TracePrefix = "[CMD] "
		Logger = logger.GetInstance(argparser.Command.Name, log.Ldate|log.Ltime|log.Lshortfile)
		fallthrough
	case len(opts.Verbose) >= 1:
		logger.Verbose = true
		shell.VerboseFunc = func(c *shell.Command) {
			Logger.Command(c.ToString())
		}
		fallthrough
	default:
		if Logger == nil {
			Logger = logger.GetInstance(argparser.Command.Name, 0)
		}
	}

	if argCommand == "help" {
		argparser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	Logger.Main("Initialisation")
	configFile := findConfigFile()
	if configFile == "" {
		Logger.FatalExit(2, "Unable to find configuration file (searched  %s)", strings.Join(validConfigFiles, " "))
	}
	Logger.Step("found configuration file %s", configFile)

	sync.Logger = Logger
	config := sync.NewConfigParser(configFile)

	switch argCommand {
	case "show":
		//----------------------
		// Show
		//----------------------
		config.ShowConfiguration()
	case "sync":
		//----------------------
		// Sync
		//----------------------
		server := getArgServer(config, "sync")
		confServer, err := config.GetSyncServer(server)
		if err != nil {
			Logger.FatalErrorExit(3, err)
		}
		Logger.Step("using Server[%s]", server)
		Logger.Step("using %s", confServer.Connection.String())

		// --dump
		if opts.Dump {
			fmt.Println()
			fmt.Println(confServer.AsYaml())
		} else {
			confServer.Sync()
			Logger.Println("-> finished")
		}
	case "deploy":
		//----------------------
		// Deploy
		//----------------------
		server := getArgServer(config, "deploy")
		confServer, err := config.GetDeployServer(server)
		if err != nil {
			Logger.FatalErrorExit(3, err)
		}
		Logger.Step("using %s server", server)
		Logger.Step("using connection %s", confServer.Connection.String())
		// --dump
		if opts.Dump {
			fmt.Println()
			fmt.Println(confServer.AsYaml())
		} else {
			confServer.Deploy()
			Logger.Println("-> finished")
		}
	}

	os.Exit(0)
}
