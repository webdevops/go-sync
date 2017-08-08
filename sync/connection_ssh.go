package sync

import (
	"fmt"
)

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
