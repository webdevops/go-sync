package sync

import (
	"strings"
	"github.com/webdevops/go-shell"
)

func (database *database) localSshDump(additionalArgs []string, useFilter bool) []interface{} {
	var args []string

	if database.Local.User != "" {
		args = append(args, "-u" + database.Local.User)
	}

	if database.Local.Password != "" {
		args = append(args, "-p" + database.Local.Password)
	}

	if database.Local.Hostname != "" {
		args = append(args, "-h" + database.Local.Hostname)
	}

	if database.Local.Port != "" {
		args = append(args, "-P" + database.Local.Port)
	}

	if len(args) > 0 {
		args = append(args, additionalArgs...)
	}

	// exclude
	excludeArgs, includeArgs := database.mysqlTableFilter(&database.Local.Connection, "local");
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// schema
	args = append(args, database.Local.Schema)

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	return database.Local.Connection.RemoteCommandBuilder("mysqldump", args...)
}

func (database *database) localCommandInterface(command string, args ...string) []interface{} {
	var ret []interface{}

	if database.Local.Connection.Type == "" {
		database.Local.Connection.Type = "local"

		// autodetection
		if database.Local.Connection.Docker != "" {
			database.Local.Connection.Type = "docker"
		}

		if database.Local.Connection.Hostname != "" {
			database.Local.Connection.Type = "ssh"
		}
	}

	switch database.Local.Connection.Type {
	case "local":
		ret = ShellCommandInterfaceBuilder(command, args...)
	case "ssh":
		ret = database.Local.Connection.RemoteCommandBuilder(command, args...)
	}

	return ret
}

func (database *database) localMysqlCmdBuilder(args ...string) []interface{} {
	args = append(args, "-BN")

	if database.Local.User != "" {
		args = append(args, "-u" + database.Local.User)
	}

	if database.Local.Password != "" {
		args = append(args, "-p" + database.Local.Password)
	}

	if database.Local.Hostname != "" {
		args = append(args, "-h" + database.Local.Hostname)
	}

	if database.Local.Port != "" {
		args = append(args, "-P" + database.Local.Port)
	}

	if database.Local.Schema != "" {
		args = append(args, database.Local.Schema)
	}

	return database.Local.Connection.RemoteCommandBuilder("mysql", args...)
}

func (database *database) localMysqlTableList() []string {
	sqlStmt := "SHOW TABLES"
	cmd := shell.Cmd("echo", sqlStmt).Pipe(database.localMysqlCmdBuilder()...)
	output := cmd.Run().Stdout.String()

	outputString := strings.TrimSpace(string(output))
	ret := strings.Split(outputString, "\n")

	return ret
}
