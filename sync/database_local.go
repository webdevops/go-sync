package sync

import "github.com/webdevops/go-shell"

func (database *Database) localMysqldumpCmdBuilder(additionalArgs []string, useFilter bool) []interface{} {
	var args []string


	if database.Local.User != "" {
		args = append(args, shell.Quote("-u" + database.Local.User))
	}

	if database.Local.Password != "" {
		args = append(args, shell.Quote("-p" + database.Local.Password))
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
	excludeArgs, includeArgs := database.mysqlTableFilter(&database.Local.Connection, "local");
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// add custom options (raw)
	if database.Local.Options.Mysqldump != "" {
		args = append(args, database.Local.Options.Mysqldump)
	}

	// schema
	args = append(args, shell.Quote(database.Local.Schema))

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	return database.Local.Connection.RawCommandBuilder("mysqldump", args...)
}

func (database *Database) localMysqlCmdBuilder(args ...string) []interface{} {
	args = append(args, "-BN")

	if database.Local.User != "" {
		args = append(args, shell.Quote("-u" + database.Local.User))
	}

	if database.Local.Password != "" {
		args = append(args, shell.Quote("-p" + database.Local.Password))
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

	if database.Local.Schema != "" {
		args = append(args, shell.Quote(database.Local.Schema))
	}

	return database.Local.Connection.RawCommandBuilder("mysql", args...)
}

