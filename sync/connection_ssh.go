package sync

import (
	"fmt"
	"github.com/webdevops/go-shell"
	"strings"
)

func (connection *Connection) SshCompressedCommandBuilder(command string, args ...string) []interface{} {
	originalCmd := []string{
		shell.Quote(command),
	}

	for _, val := range args {
		originalCmd = append(originalCmd, shell.Quote(val))
	}

	inlineCommand := fmt.Sprintf("%s | gzip --stdout", strings.Join(originalCmd, " "))

	return connection.SshCommandBuilder(inlineCommand)
}

func (connection *Connection) SshCommandBuilder(command string, args ...string) []interface{} {
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

func (connection *Connection) SshConnectionHostnameString() string {
	if connection.User != "" {
		return fmt.Sprintf("%s@%s", connection.User, connection.Hostname)
	} else {
		return connection.Hostname
	}
}
