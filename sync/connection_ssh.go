package sync

import (
	"fmt"
	"strings"
	"github.com/webdevops/go-shell"
)

func (connection *Connection) SshCommandBuilder(command string, args ...string) []interface{} {
	remoteCmdParts := []string{command}
	for _, val := range args {
		remoteCmdParts = append(remoteCmdParts, val)
	}
	remoteCmd := shell.Quote(strings.Join(remoteCmdParts, " "))

	sshArgs := []string{
		"-oBatchMode=yes",
		"-oPasswordAuthentication=no",
		connection.SshConnectionHostnameString(),
		"--",
		remoteCmd,
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
