package sync

import (
	"github.com/webdevops/go-shell"
)

func (database *DatabaseMysql) localMysqldumpCmdBuilder(additionalArgs []string, useFilter bool) []interface{} {
	var args []string

	connection := database.Local.Connection.GetInstance().Clone()

	if database.Local.User != "" {
		args = append(args, shell.Quote("-u" + database.Local.User))
	}

	if database.Local.Password != "" {
		connection.Environment.Set("MYSQL_PWD", database.Local.Password)
	}

	if database.Local.Hostname != "" {
		args = append(args, shell.Quote("-h" + database.Local.Hostname))
	}

	if database.Local.Port != "" {
		args = append(args, shell.Quote("-P" + database.Local.Port))
	}

	if len(args) > 0 {
		args = append(args, additionalArgs...)
	}

	// exclude
	excludeArgs, includeArgs := database.tableFilter("local");
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// add custom options (raw)
	if database.Local.Options.Mysqldump != "" {
		args = append(args, database.Local.Options.Mysqldump)
	}

	// database
	args = append(args, shell.Quote(database.Local.Db))

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	return connection.RawCommandBuilder("mysqldump", args...)
}

func (database *DatabaseMysql) localMysqlCmdBuilder(args ...string) []interface{} {
	connection := database.Local.Connection.GetInstance().Clone()

	args = append(args, "-BN")

	if database.Local.User != "" {
		args = append(args, shell.Quote("-u" + database.Local.User))
	}

	if database.Local.Password != "" {
		connection.Environment.Set("MYSQL_PWD", database.Local.Password)
	}

	if database.Local.Hostname != "" {
		args = append(args, shell.Quote("-h" + database.Local.Hostname))
	}

	if database.Local.Port != "" {
		args = append(args, shell.Quote("-P" + database.Local.Port))
	}

	// add custom options (raw)
	if database.Local.Options.Mysql != "" {
		args = append(args, database.Local.Options.Mysql)
	}

	if database.Local.Db != "" {
		args = append(args, shell.Quote(database.Local.Db))
	}

	return connection.RawCommandBuilder("mysql", args...)
}

