package sync

import (
	"fmt"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
)

type DatabaseMysql struct {
	Database
}

func (database *DatabaseMysql) tableFilter(connection *commandbuilder.Connection, connectionType string) ([]string, []string) {
	var exclude []string
	var include []string

	var tableList []string

	if (connectionType == "local") {
		if len(database.cacheLocalTableList) == 0 {
			Logger.Step("get list of mysql tables for table filter")
			database.cacheLocalTableList = database.tableList(connectionType)
		}

		tableList = database.cacheLocalTableList
	} else {
		if len(database.cacheRemoteTableList) == 0 {
			Logger.Step("get list of mysql tables for table filter")
			database.cacheRemoteTableList = database.tableList(connectionType)
		}

		tableList = database.cacheRemoteTableList
	}


	// calc excludes
	excludeTableList := database.Filter.CalcExcludes(tableList)
	for _, table := range excludeTableList {
		exclude  = append(exclude, shell.Quote(fmt.Sprintf("--ignore-table=%s.%s", database.Schema, table)))
	}

	// calc includes
	includeTableList := database.Filter.CalcIncludes(tableList)
	for _, table := range includeTableList {
		include  = append(include, table)
	}

	return exclude, include
}

func (database *DatabaseMysql) mysqlCommandBuilder(direction string, args ...string) []interface{} {
	if direction == "local" {
		return database.localMysqlCmdBuilder(args...)
	} else {
		return database.remoteMysqlCmdBuilder(args...)
	}
}

func (database *DatabaseMysql) tableList(connectionType string) []string {
	sqlStmt := "SHOW TABLES"

	cmd := shell.Cmd("echo", shell.Quote(sqlStmt)).Pipe(database.mysqlCommandBuilder(connectionType)...)
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

func (database *DatabaseMysql) quote(value string) string {
	return "'" + strings.Replace(value, "'", "\\'", -1) + "'"
}

func (database *DatabaseMysql) quoteIdentifier(value string) string {
	return "`" + strings.Replace(value, "`", "\\`", -1) + "`"
}
