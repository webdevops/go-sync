package sync

import (
	"strings"
	"github.com/webdevops/go-shell"
	"fmt"
)

func (execution *Execution) String(server *Server) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Type:%s", execution.GetType()))


	if execution.Workdir != "" {
		parts = append(parts, fmt.Sprintf("Workdir:%s", execution.Workdir))
	}
	parts = append(parts, fmt.Sprintf("Command:%s", execution.Command))

	return fmt.Sprintf("Exec[%s]", strings.Join(parts[:]," "))
}

func (execution *Execution) Execute(server *Server) {
	shell.Cmd(execution.commandBuilder(server)...)
}

func (execution *Execution) commandBuilder(server *Server) []interface{} {
	var connection Connection

	switch execution.GetType() {
	case "local":
		connection = Connection{Type:"local"}
	case "remote":
		connection = server.Connection
	}

	if execution.Workdir != "" {
		connection.WorkDir = execution.Workdir
	}

	return connection.CommandBuilder(execution.Command)
}

func (execution *Execution) GetType() string {
	var ret string

	switch strings.ToLower(execution.Type) {
	case "":
		fallthrough
	case "local":
		ret = "local"
	case "remote":
		ret = "remote"
	default:
		panic(execution)
	}

	return ret
}
