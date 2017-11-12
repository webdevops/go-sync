package sync

import "github.com/webdevops/go-shell"

func (database *DatabasePostgres) localPgdumpCmdBuilder(additionalArgs []string, useFilter bool) []interface{} {
	var args []string

	connection := database.Local.Connection.GetInstance().Clone()

	// add custom options (raw)
	if database.Local.Options.Pgdump != nil {
		args = append(args, database.Local.Options.Pgdump.Array()...)
	}

	if database.Local.User != "" {
		args = append(args, "-U", shell.Quote(database.Local.User))
	}

	if database.Local.Password != "" {
		connection.Environment.Set("PGPASSWORD", database.Local.Password)
	}

	if database.Local.Hostname != "" {
		args = append(args, "-h", shell.Quote(database.Local.Hostname))
	}

	if database.Local.Port != "" {
		args = append(args, "-p", shell.Quote(database.Local.Port))
	}

	if len(args) > 0 {
		args = append(args, additionalArgs...)
	}

	// exclude
	excludeArgs, includeArgs := database.tableFilter("local");
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// database
	args = append(args, shell.Quote(database.Local.Db))

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	return connection.RawCommandBuilder("pg_dump", args...)
}

func (database *DatabasePostgres) localPsqlCmdBuilder(additonalArgs ...string) []interface{} {
	var args []string

	connection := database.Local.Connection.GetInstance().Clone()

	// add custom options (raw)
	if database.Local.Options.Psql != nil {
		args = append(args, database.Local.Options.Psql.Array()...)
	}

	args = append(args, "-t")

	if database.Local.User != "" {
		args = append(args, "-U", shell.Quote(database.Local.User))
	}

	if database.Local.Password != "" {
		connection.Environment.Set("PGPASSWORD", database.Local.Password)
	}

	if database.Local.Hostname != "" {
		args = append(args, "-h", shell.Quote(database.Local.Hostname))
	}

	if database.Local.Port != "" {
		args = append(args, "-p", shell.Quote(database.Local.Port))
	}

	if len(additonalArgs) > 0 {
		args = append(args, additonalArgs...)
	}

	if database.Local.Db != "" {
		args = append(args, shell.Quote(database.Local.Db))
	}

	return connection.RawCommandBuilder("psql", args...)
}
