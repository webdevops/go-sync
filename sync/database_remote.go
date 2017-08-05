package sync

import (
	"strings"
	"github.com/webdevops/go-shell"
)

func (database *database) remoteSshDump(additionalArgs []string, useFilter bool) []interface{} {
	var args []string

	if database.User != "" {
		args = append(args, "-u" + database.User)
	}

	if database.Password != "" {
		args = append(args, "-p" + database.Password)
	}

	if database.Hostname != "" {
		args = append(args, "-h" + database.Hostname)
	}

	if database.Port != "" {
		args = append(args, "-P" + database.Port)
	}

	if len(args) > 0 {
		args = append(args, additionalArgs...)
	}

	// exclude
	excludeArgs, includeArgs := database.mysqlTableFilter(&database.remoteConnection, "remote");
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// schema
	args = append(args, database.Schema)

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	return database.remoteConnection.SshCompressedCommandBuilder("mysqldump", args...)
}

func (database *database) remoteMysqlCmdBuilder(args ...string) []interface{} {
	args = append(args, "-BN")

	if database.User != "" {
		args = append(args, "-u" + database.User)
	}

	if database.Password != "" {
		args = append(args, "-p" + database.Password)
	}

	if database.Hostname != "" {
		args = append(args, "-h" + database.Hostname)
	}

	if database.Port != "" {
		args = append(args, "-P" + database.Port)
	}

	if database.Schema != "" {
		args = append(args, database.Schema)
	}

	return database.remoteConnection.RemoteCommandBuilder("mysql", args...)
}

func (database *database) remoteMysqlTableList() []string {
	sqlStmt := "SHOW TABLES"
	cmd := shell.Cmd("echo", sqlStmt).Pipe(database.remoteMysqlCmdBuilder()...)
	output := cmd.Run().Stdout.String()

	outputString := strings.TrimSpace(string(output))
	ret := strings.Split(outputString, "\n")

	return ret
}
