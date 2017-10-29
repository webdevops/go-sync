package sync

import (
	"strings"
	"github.com/webdevops/go-shell"
	"fmt"
	"github.com/webdevops/go-shell/commandbuilder"
)

func (execution *Execution) String(server *Server) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Type:%s", execution.GetType()))


	if execution.Workdir != "" {
		parts = append(parts, fmt.Sprintf("Workdir:%s", execution.Workdir))
	}
	parts = append(parts, fmt.Sprintf("Command:%s", execution.Command.ToString(" ")))

	return fmt.Sprintf("Exec[%s]", strings.Join(parts[:]," "))
}

func (execution *Execution) Execute(server *Server) {
	cmd := execution.commandBuilder(server)
	run := shell.Cmd(cmd...).Run()

	Logger.Verbose(run.Stdout.String())
}

// Create commandBuilder for execution
func (execution *Execution) commandBuilder(server *Server) []interface{} {
	var connection commandbuilder.Connection

	switch execution.GetType() {
	case "local":
		connection = commandbuilder.Connection{Type:"local"}
	case "remote":
		connection = *(server.Connection.GetInstance().Clone())
	}

	// set working directory
	if execution.Workdir != "" {
		connection.Workdir = execution.Workdir
	}

	// set environment
	connection.Environment.Clear()
	for _, val := range execution.Environment {
		connection.Environment.Set(val.Name, val.Value)
	}

	if len(execution.Command.Multi) >= 1 {
		// multi element command (use safer quoting)
		return connection.ShellCommandBuilder(execution.Command.Multi...)
	} else {
		// single string command (use as is)
		return connection.RawShellCommandBuilder(execution.Command.Single)
	}
}

// Get execution type (local or remote)
func (execution *Execution) GetType() (execType string) {
	switch strings.ToLower(execution.Type) {
	case "":
		fallthrough
	case "local":
		execType = "local"
	case "remote":
		execType = "remote"
	default:
		panic(execution)
	}

	return
}
