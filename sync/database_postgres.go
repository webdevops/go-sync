package sync

import (
	"fmt"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
)

type DatabasePostgres struct {
	Database
}

func (database *DatabasePostgres) tableFilter(connection *commandbuilder.Connection, connectionType string) ([]string, []string) {
	var exclude []string
	var include []string

	var tableList []string

	if (connectionType == "local") {
		if len(database.cacheLocalTableList) == 0 {
			Logger.Step("get list of postgres tables for table filter")
			database.cacheLocalTableList = database.tableList(connectionType)
		}

		tableList = database.cacheLocalTableList
	} else {
		if len(database.cacheRemoteTableList) == 0 {
			Logger.Step("get list of postgres tables for table filter")
			database.cacheRemoteTableList = database.tableList(connectionType)
		}

		tableList = database.cacheRemoteTableList
	}

	// calc excludes
	excludeTableList := database.Filter.CalcExcludes(tableList)
	for _, table := range excludeTableList {
		exclude  = append(exclude, shell.Quote(fmt.Sprintf("--exclude-table=%s", table)))
	}

	// calc includes
	includeTableList := database.Filter.CalcIncludes(tableList)
	for _, table := range includeTableList {
		include  = append(include, shell.Quote(fmt.Sprintf("--table=%s", table)))
	}
	
	return exclude, include
}

func (database *DatabasePostgres) psqlCommandBuilder(direction string, args ...string) []interface{} {
	if direction == "local" {
		return database.localPsqlCmdBuilder(args...)
	} else {
		return database.remotePsqlCmdBuilder(args...)
	}
}

func (database *DatabasePostgres) tableList(connectionType string) []string {
	sqlStmt := `SELECT table_name
                  FROM information_schema.tables
                 WHERE table_type = 'BASE TABLE'
                   AND table_catalog = %s`
	sqlStmt = fmt.Sprintf(sqlStmt, database.quote(database.Schema))

	cmd := shell.Cmd("echo", shell.Quote(sqlStmt)).Pipe(database.psqlCommandBuilder(connectionType)...)
	output := cmd.Run().Stdout.String()

	outputString := strings.TrimSpace(string(output))
	tmp := strings.Split(outputString, "\n")

	// trim spaces from tables
	ret := make([]string, len(tmp))
	for _, table := range tmp {
		ret = append(ret, strings.TrimSpace(table))
	}

	return ret
}

func (database *DatabasePostgres) quote(value string) string {
	return "'" + strings.Replace(value, "'", "\\'", -1) + "'"
}

func (database *DatabasePostgres) quoteIdentifier(value string) string {
	return "\"" + strings.Replace(value, "\"", "\\\"", -1) + "\""
}