package sync

import (
	"strings"
	"github.com/webdevops/go-shell"
	"fmt"
)

var (
	ConnectionSshArguments = []string{"-oBatchMode=yes -oPasswordAuthentication=no"}
	ConnectionDockerArguments = []string{"exec", "-i"}
)

func (connection *Connection) CommandBuilder(command string, args ...string) []interface{} {
	//args = shell.QuoteValues(args...)

	return connection.RawCommandBuilder(command, args...)
}

func (connection *Connection) RawCommandBuilder(command string, args ...string) []interface{} {
	var ret []interface{}

	// if workdir is set
	// use shell'ed command builder
	if connection.WorkDir != "" {
		shellArgs := []string{command}
		shellArgs = append(shellArgs, args...)
		return connection.RawShellCommandBuilder(shellArgs...)
	}

	switch connection.GetType() {
	case "local":
		ret = connection.LocalCommandBuilder(command, args...)
	case "ssh":
		ret = connection.SshCommandBuilder(command, args...)
	case "ssh+docker":
		fallthrough
	case "docker":
		ret = connection.DockerCommandBuilder(command, args...)
	default:
		panic(connection)
	}

	return ret
}

func (connection *Connection) ShellCommandBuilder(args ...string) []interface{} {
	args = shell.QuoteValues(args...)
	return connection.RawShellCommandBuilder(args...)
}

func (connection *Connection) RawShellCommandBuilder(args ...string) []interface{} {
	var ret []interface{}

	inlineArgs := []string{}

	for _, val := range args {
		inlineArgs = append(inlineArgs, val)
	}

	inlineCommand := strings.Join(inlineArgs, " ")

	if connection.WorkDir != "" {
		// prepend cd in front of command to change work dir
		inlineCommand = fmt.Sprintf("cd %s ; %s", shell.Quote(connection.WorkDir), inlineCommand)
	}

	inlineCommand = shell.Quote(inlineCommand)

	switch connection.GetType() {
	case "local":
		ret = connection.LocalCommandBuilder("/bin/sh", "-c", inlineCommand)
	case "ssh":
		ret = connection.SshCommandBuilder("/bin/sh", "-c", inlineCommand)
	case "ssh+docker":
		fallthrough
	case "docker":
		ret = connection.DockerCommandBuilder("/bin/sh", "-c", inlineCommand)
	default:
		panic(connection)
	}

	return ret
}

func (connection *Connection) GetType() string {
	var connType string

	// autodetection
	if (connection.Type == "") || (connection.Type == "auto") {
		connection.Type = "local"

		if (connection.Docker != "") && connection.Hostname != "" {
			connection.Type = "ssh+docker"
		} else if connection.Docker != "" {
			connection.Type = "docker"
		} else if connection.Hostname != "" {
			connection.Type = "ssh"
		}
	}

	switch connection.Type {
	case "local":
		connType = "local"
	case "ssh":
		connType = "ssh"
	case "docker":
		connType = "docker"
	case "ssh+docker":
		connType = "ssh+docker"
	default:
		Logger.FatalExit(1, "Unknown connection type \"%s\"", connType)
	}

	return connType
}

func (connection *Connection) String() string {
	var parts []string

	connType := connection.GetType()
	parts = append(parts, fmt.Sprintf("Type:%s", connType))

	switch connType {
	case "ssh":
		parts = append(parts, fmt.Sprintf("SSH:%s", connection.SshConnectionHostnameString()))
	case "docker":
		parts = append(parts, fmt.Sprintf("Docker:%s", connection.Docker))
	case "ssh+docker":
		parts = append(parts, fmt.Sprintf("SSH:%s", connection.SshConnectionHostnameString()))
		parts = append(parts, fmt.Sprintf("Docker:%s", connection.Docker))
	default:
	}

	return fmt.Sprintf("Exec[%s]", strings.Join(parts[:]," "))
}
