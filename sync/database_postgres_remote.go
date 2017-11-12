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
		connection.Environment.Set("PGPASSWORD", database.Password)
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

	// database
	args = append(args, shell.Quote(database.Db))

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	cmd := []string{"pg_dump"}

	// add custom options (raw)
	if database.Options.Pgdump != nil {
		cmd = append(cmd, database.Options.Pgdump.Array()...)
	}

	cmd = append(cmd, args...)
	cmd = append(cmd, "|", "gzip", "--stdout")

	return connection.RawShellCommandBuilder(cmd...)
}

func (database *DatabasePostgres) remotePsqlCmdBuilder(additonalArgs ...string) []interface{} {
	var args []string

	connection := database.Connection.GetInstance().Clone()

	// append options in raw
	if database.Options.Psql != nil {
		args = append(args, database.Options.Psql.Array()...)
	}

	args = append(args, "-t")

	if database.User != "" {
		args = append(args, "-U", shell.Quote(database.User))
	}

	if database.Password != "" {
		connection.Environment.Set("PGPASSWORD", database.Password)
	}

	if database.Hostname != "" {
		args = append(args, "-h", shell.Quote(database.Hostname))
	}

	if database.Port != "" {
		args = append(args, "-p", shell.Quote(database.Port))
	}

	if database.Db != "" {
		args = append(args, shell.Quote(database.Db))
	}

	if len(additonalArgs) > 0 {
		args = append(args, additonalArgs...)
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
		connection.Environment.Set("MYSQL_PWD", database.Password)
	}

	if database.Hostname != "" {
		args = append(args, "-h", shell.Quote(database.Hostname))
	}

	if database.Port != "" {
		args = append(args, "-p", shell.Quote(database.Port))
	}

	// add custom options (raw)
	if database.Options.Psql != nil {
		args = append(args, database.Options.Psql.Array()...)
	}

	if database.Db != "" {
		args = append(args, shell.Quote(database.Db))
	}

	cmd := []string{"gunzip", "--stdout", "|", "psql", strings.Join(args, " ")}

	return connection.RawShellCommandBuilder(cmd...)
}
