package main

import (
	"os"
	"log"
	"fmt"
	"runtime/debug"
	flags "github.com/jessevdk/go-flags"
	"github.com/webdevops/go-shell"
	"./logger"
)

const (
	// application informations
	Name    = "gosync"
	Author  = "webdevops.io"
	Version = "0.4.1"

	// self update informations
	GithubOrganization  = "webdevops"
	GithubRepository    = "go-sync"
	GithubAssetTemplate = "gosync-%OS%-%ARCH%"
)

var (
	Logger *logger.SyncLogger
	argparser *flags.Parser
	args []string
)

var opts struct {
	Verbose            []bool   `short:"v"  long:"verbose"                       description:"verbose mode"`
}

var validConfigFiles = []string{
	"gosync.yml",
	"gosync.yaml",
	".gosync.yml",
	".gosync.yaml",
}

func handleArgParser() {
	var err error
	argparser = flags.NewParser(&opts, flags.Default)
	argparser.CommandHandler = func(command flags.Commander, args []string) error {
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

		return command.Execute(args)
	}

	argparser.AddCommand("version", "Show version", fmt.Sprintf("Show %s version", Name), &VersionCommand{Name:Name, Version:Version, Author:Author})
	argparser.AddCommand("self-update", "Self update", "Run self update of this application", &SelfUpdateCommand{GithubOrganization:GithubOrganization, GithubRepository:GithubRepository, GithubAssetTemplate:GithubAssetTemplate, CurrentVersion:Version})

	argparser.AddCommand("list", "List server configurations", "List server configurations", &ListCommand{})
	argparser.AddCommand("sync", "Sync from server", "Sync filesystem and databases from server", &SyncCommand{})
	argparser.AddCommand("deploy", "Deploy to server", "Deploy filesystem and databases to server", &DeployCommand{})

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

	shell.SetDefaultShell("bash")
	handleArgParser()

	os.Exit(0)
}
