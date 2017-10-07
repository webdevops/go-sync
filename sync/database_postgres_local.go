package sync

import "github.com/webdevops/go-shell"

func (database *DatabasePostgres) localPgdumpCmdBuilder(additionalArgs []string, useFilter bool) []interface{} {
	connection := database.Local.Connection.Clone()
	var args []string

	if database.Local.User != "" {
		args = append(args, "-U", shell.Quote(database.Local.User))
	}

	if database.Local.Password != "" {
		connection.Environment["PGPASSWORD"] = database.Local.Password
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
	excludeArgs, includeArgs := database.tableFilter(&database.Local.Connection, "local");
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// add custom options (raw)
	if database.Local.Options.Pgdump != "" {
		args = append(args, database.Local.Options.Pgdump)
	}

	// schema
	args = append(args, shell.Quote(database.Local.Schema))

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	return connection.RawCommandBuilder("pg_dump", args...)
}

func (database *DatabasePostgres) localPsqlCmdBuilder(args ...string) []interface{} {
	connection := database.Local.Connection.Clone()
	args = append(args, "-t")

	if database.Local.User != "" {
		args = append(args, "-U", shell.Quote(database.Local.User))
	}

	if database.Local.Password != "" {
		connection.Environment["PGPASSWORD"] = database.Local.Password
	}

	if database.Local.Hostname != "" {
		args = append(args, "-h", shell.Quote(database.Local.Hostname))
	}

	if database.Local.Port != "" {
		args = append(args, "-p", shell.Quote(database.Local.Port))
	}

	// add custom options (raw)
	if database.Local.Options.Psql != "" {
		args = append(args, database.Local.Options.Psql)
	}

	if database.Local.Schema != "" {
		args = append(args, shell.Quote(database.Local.Schema))
	}

	return connection.RawCommandBuilder("psql", args...)
}
