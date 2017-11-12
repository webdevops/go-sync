package sync

import (
	"strings"
	"github.com/webdevops/go-shell"
)

func (database *DatabaseMysql) remoteMysqldumpCmdBuilder(additionalArgs []string, useFilter bool) []interface{} {
	var args []string

	connection := database.Connection.GetInstance().Clone()

	if database.User != "" {
		args = append(args, shell.Quote("-u" + database.User))
	}

	if database.Password != "" {
		connection.Environment.Set("MYSQL_PWD", database.Password)
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
	excludeArgs, includeArgs := database.tableFilter("remote");
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// database
	args = append(args, shell.Quote(database.Db))

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	cmd := []string{"mysqldump"}

	// add custom options (raw)
	if database.Options.Mysqldump != nil {
		cmd = append(cmd, database.Options.Mysqldump.Array()...)
	}

	cmd = append(cmd, args...)
	cmd = append(cmd, "|", "gzip", "--stdout")

	return connection.RawShellCommandBuilder(cmd...)
}

func (database *DatabaseMysql) remoteMysqlCmdBuilder(additonalArgs ...string) []interface{} {
	var args []string

	connection := database.Connection.GetInstance().Clone()

	// append options in raw
	if database.Options.Mysql != nil {
		args = append(args, database.Options.Mysql.Array()...)
	}

	args = append(args, "-BN")

	if database.User != "" {
		args = append(args, shell.Quote("-u" + database.User))
	}

	if database.Password != "" {
		connection.Environment.Set("MYSQL_PWD", database.Password)
	}

	if database.Hostname != "" {
		args = append(args, shell.Quote("-h" + database.Hostname))
	}

	if database.Port != "" {
		args = append(args, shell.Quote("-P" + database.Port))
	}

	if database.Db != "" {
		args = append(args, shell.Quote(database.Db))
	}

	if len(additonalArgs) > 0 {
		args = append(args, additonalArgs...)
	}

	return connection.RawCommandBuilder("mysql", args...)
}


func (database *DatabaseMysql) remoteMysqlCmdBuilderUncompress(additonalArgs ...string) []interface{} {
	var args []string

	connection := database.Connection.GetInstance().Clone()

	// add custom options (raw)
	if database.Options.Mysql != nil {
		args = append(args, database.Options.Mysql.Array()...)
	}

	args = append(args, "-BN")

	if database.User != "" {
		args = append(args, shell.Quote("-u" + database.User))
	}

	if database.Password != "" {
		connection.Environment.Set("MYSQL_PWD", database.Password)
	}

	if database.Hostname != "" {
		args = append(args, shell.Quote("-h" + database.Hostname))
	}

	if database.Port != "" {
		args = append(args, shell.Quote("-P" + database.Port))
	}

	if len(additonalArgs) > 0 {
		args = append(args, additonalArgs...)
	}

	if database.Db != "" {
		args = append(args, shell.Quote(database.Db))
	}

	cmd := []string{"gunzip", "--stdout", "|", "mysql", strings.Join(args, " ")}

	return connection.RawShellCommandBuilder(cmd...)
}
