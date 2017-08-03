package sync

import (
	"fmt"
	"github.com/webdevops/go-shell"
	"strings"
)

type connection struct {
	Type string
	Hostname string
	User string
	Password string
	Docker string
	Compression string
}

func (connection *connection) SshCompressedCommandBuilder(command string, args ...string) []interface{} {
	originalCmd := []string{
		shell.Quote(command),
	}

	for _, val := range args {
		originalCmd = append(originalCmd, shell.Quote(val))
	}

	inlineCommand := fmt.Sprintf("%s | gzip --stdout", strings.Join(originalCmd, " "))

	return connection.sshCommandBuilder(inlineCommand)
}

func (connection *connection) RemoteCommandBuilder(command string, args ...string) []interface{} {
	var ret []interface{}

	if connection.Type == "" {
		connection.Type = "local"

		// autodetection
		if connection.Docker != "" {
			connection.Type = "docker"
		}

		if connection.Hostname != "" {
			connection.Type = "ssh"
		}
	}

	switch connection.GetType() {
	case "local":
		ret = ShellCommandInterfaceBuilder(command, args...)
	case "ssh":
		ret = connection.sshCommandBuilder(command, args...)
	case "docker":
		ret = connection.dockerCommandBuilder(command, args...)
	}

	return ret
}

func (connection *connection) sshCommandBuilder(command string, args ...string) []interface{} {
	sshArgs := []string{
		"-oBatchMode=yes",
		"-oPasswordAuthentication=no",
		connection.SshConnectionHostnameString(),
		"--",
		command,
	}

	for _, val := range args {
		sshArgs = append(sshArgs, val)
	}

	return ShellCommandInterfaceBuilder("ssh", sshArgs...)
}

func (connection *connection) dockerCommandBuilder(cmd string, args ...string) []interface{} {
	dockerArgs := []string{
		"exec",
		"-i",
		DockerGetContainerId(connection.Docker),
		cmd,
	}

	for _, val := range args {
		dockerArgs = append(dockerArgs, val)
	}

	return ShellCommandInterfaceBuilder("docker", dockerArgs...)
}

func (connection *connection) SshConnectionHostnameString() string {
	if connection.User != "" {
		return fmt.Sprintf("%s@%s", connection.User, connection.Hostname)
	} else {
		return connection.Hostname
	}
}

func (connection *connection) GetType() string {
	var connType string

	switch connection.Type {
	case "":
		fallthrough
	case "ssh":
		connType = "ssh"
	case "docker":
		connType = "docker"
	default:
		Logger.FatalExit(1, "Unknown connection type \"%s\"", connType)
	}

	return connType
}
