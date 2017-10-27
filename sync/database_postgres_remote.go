package sync

import (
	"strings"
	"github.com/webdevops/go-shell"
)

func (database *DatabasePostgres) remotePgdumpCmdBuilder(additionalArgs []string, useFilter bool) []interface{} {
	connection := database.Connection.GetInstance().Clone()
	var args []string

	if database.User != "" {
		args = append(args, "-U", shell.Quote(database.User))
	}

	if database.Password != "" {
		connection.Environment["PGPASSWORD"] = database.Password
	}

	if database.Hostname != "" {
		args = append(args, "-h", shell.Quote(database.Hostname))
	}

	if database.Port != "" {
		args = append(args, "-p", shell.Quote(database.Port))
	}

	if len(args) > 0 {
		args = append(args, additionalArgs...)
	}

	// exclude
	excludeArgs, includeArgs := database.tableFilter("remote")
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// schema
	args = append(args, shell.Quote(database.Schema))

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	cmd := []string{"pg_dump"}

	// add custom options (raw)
	if database.Options.Pgdump != "" {
		cmd = append(cmd, database.Options.Pgdump)
	}

	cmd = append(cmd, args...)
	cmd = append(cmd, "|", "gzip", "--stdout")

	return connection.RawShellCommandBuilder(cmd...)
}

func (database *DatabasePostgres) remotePsqlCmdBuilder(args ...string) []interface{} {
	connection := database.Connection.GetInstance().Clone()
	args = append(args, "-t")

	if database.User != "" {
		args = append(args, "-U", shell.Quote(database.User))
	}

	if database.Password != "" {
		connection.Environment["PGPASSWORD"] = database.Password
	}

	if database.Hostname != "" {
		args = append(args, "-h", shell.Quote(database.Hostname))
	}

	if database.Port != "" {
		args = append(args, "-p", shell.Quote(database.Port))
	}

	if database.Schema != "" {
		args = append(args, shell.Quote(database.Schema))
	}

	// append options in raw
	if database.Options.Psql != "" {
		args = append(args, database.Options.Psql)
	}

	return connection.RawCommandBuilder("psql", args...)
}


func (database *DatabasePostgres) remotePsqlCmdBuilderUncompress(args ...string) []interface{} {
	connection := database.Connection.GetInstance().Clone()
	args = append(args, "-t")

	if database.User != "" {
		args = append(args, "-U", shell.Quote(database.User))
	}

	if database.Password != "" {
		connection.Environment["MYSQL_PWD"] = database.Password
	}

	if database.Hostname != "" {
		args = append(args, "-h", shell.Quote(database.Hostname))
	}

	if database.Port != "" {
		args = append(args, "-p", shell.Quote(database.Port))
	}

	// add custom options (raw)
	if database.Options.Psql != "" {
		args = append(args, database.Options.Psql)
	}

	if database.Schema != "" {
		args = append(args, shell.Quote(database.Schema))
	}

	cmd := []string{"gunzip", "--stdout", "|", "psql", strings.Join(args, " ")}

	return connection.RawShellCommandBuilder(cmd...)
}
