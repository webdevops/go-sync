package sync

import (
	"strings"
	"github.com/webdevops/go-shell"
)

func (database *DatabaseMysql) remoteMysqldumpCmdBuilder(additionalArgs []string, useFilter bool) []interface{} {
	connection := database.Connection.Clone()
	var args []string

	if database.User != "" {
		args = append(args, shell.Quote("-u" + database.User))
	}

	if database.Password != "" {
		connection.Environment["MYSQL_PWD"] = database.Password
	}

	if database.Hostname != "" {
		args = append(args, shell.Quote("-h" + database.Hostname))
	}

	if database.Port != "" {
		args = append(args, shell.Quote("-P" + database.Port))
	}

	if len(args) > 0 {
		args = append(args, additionalArgs...)
	}

	// exclude
	excludeArgs, includeArgs := database.tableFilter(&database.Connection, "remote");
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// schema
	args = append(args, shell.Quote(database.Schema))

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	cmd := []string{"mysqldump"}

	// add custom options (raw)
	if database.Options.Mysqldump != "" {
		cmd = append(cmd, database.Options.Mysqldump)
	}

	cmd = append(cmd, args...)
	cmd = append(cmd, "|", "gzip", "--stdout")

	return connection.RawShellCommandBuilder(cmd...)
}

func (database *DatabaseMysql) remoteMysqlCmdBuilder(args ...string) []interface{} {
	connection := database.Connection.Clone()
	args = append(args, "-BN")

	if database.User != "" {
		args = append(args, shell.Quote("-u" + database.User))
	}

	if database.Password != "" {
		connection.Environment["MYSQL_PWD"] = database.Password
	}

	if database.Hostname != "" {
		args = append(args, shell.Quote("-h" + database.Hostname))
	}

	if database.Port != "" {
		args = append(args, shell.Quote("-P" + database.Port))
	}

	if database.Schema != "" {
		args = append(args, shell.Quote(database.Schema))
	}

	// append options in raw
	if database.Options.Mysql != "" {
		args = append(args, database.Options.Mysql)
	}

	return connection.RawCommandBuilder("mysql", args...)
}


func (database *DatabaseMysql) remoteMysqlCmdBuilderUncompress(args ...string) []interface{} {
	connection := database.Connection.Clone()
	args = append(args, "-BN")

	if database.User != "" {
		args = append(args, shell.Quote("-u" + database.User))
	}

	if database.Password != "" {
		connection.Environment["MYSQL_PWD"] = database.Password
	}

	if database.Hostname != "" {
		args = append(args, shell.Quote("-h" + database.Hostname))
	}

	if database.Port != "" {
		args = append(args, shell.Quote("-P" + database.Port))
	}

	// add custom options (raw)
	if database.Options.Mysql != "" {
		args = append(args, database.Options.Mysql)
	}

	if database.Schema != "" {
		args = append(args, shell.Quote(database.Schema))
	}

	cmd := []string{"gunzip", "--stdout", "|", "mysql", strings.Join(args, " ")}

	return connection.RawShellCommandBuilder(cmd...)
}
