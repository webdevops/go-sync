package sync

import (
	"strings"
	"github.com/webdevops/go-shell"
)

func (database *Database) remoteMysqldumpCmdBuilder(additionalArgs []string, useFilter bool) []interface{} {
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

	cmd := []string{"mysqldump"}
	cmd = append(cmd, shell.QuoteValues(args...)...)
	cmd = append(cmd, "|", "gzip", "--stdout")

	return database.remoteConnection.RawShellCommandBuilder(cmd...)
}

func (database *Database) remoteMysqlCmdBuilder(args ...string) []interface{} {
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

	return database.remoteConnection.CommandBuilder("mysql", args...)
}


func (database *Database) remoteMysqlCmdBuilderUncompress(args ...string) []interface{} {
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

	cmd := []string{"gunzip", "--stdout", "|", "mysql", strings.Join(shell.QuoteValues(args...), " ")}

	return database.remoteConnection.RawShellCommandBuilder(cmd...)
}
