package sync

import (
	"fmt"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
)


func (database *Database) mysqlTableFilter(connection *commandbuilder.Connection, connectionType string) ([]string, []string) {
	var exclude []string
	var include []string

	var tableList []string

	if (connectionType == "local") {
		if len(database.cacheLocalTableList) == 0 {
			Logger.Step("get list of mysql tables for table filter")
			database.cacheLocalTableList = database.mysqlTableList(connectionType)
		}

		tableList = database.cacheLocalTableList
	} else {
		if len(database.cacheRemoteTableList) == 0 {
			Logger.Step("get list of mysql tables for table filter")
			database.cacheRemoteTableList = database.mysqlTableList(connectionType)
		}

		tableList = database.cacheRemoteTableList
	}


	// calc excludes
	excludeTableList := database.Filter.CalcExcludes(tableList)
	for _, table := range excludeTableList {
		exclude  = append(exclude, fmt.Sprintf("--ignore-table=%s.%s", database.Schema, table))
	}

	// calc includes
	includeTableList := database.Filter.CalcIncludes(tableList)
	for _, table := range includeTableList {
		include  = append(include, table)
	}

	return exclude, include
}

func (database *Database) String(direction string) string {
	var parts, remote, local []string

	// remote
	remote = append(remote, fmt.Sprintf("Schema:%s", database.Schema))

	if database.Hostname != "" {
		remote = append(remote, fmt.Sprintf("Host:%s", database.Hostname))
	}

	if database.User != "" {
		remote = append(remote, fmt.Sprintf("User:%s", database.User))
	}

	if database.Password != "" {
		remote = append(remote, fmt.Sprintf("Passwd:%s", "*****"))
	}

	// local
	local = append(local, fmt.Sprintf("Schema:%s", database.Local.Schema))

	// build parts
	switch direction {
	case "sync":
		parts = append(parts, remote...)
		parts = append(parts, "->")
		parts = append(parts, local...)
	case "deploy":
		parts = append(parts, local...)
		parts = append(parts, "->")
		parts = append(parts, remote...)
	}

	return fmt.Sprintf("Database[%s]", strings.Join(parts[:]," "))
}

func (database *Database) mysqlCommandBuilder(direction string, args ...string) []interface{} {
	if direction == "local" {
		return database.localMysqlCmdBuilder(args...)
	} else {
		return database.remoteMysqlCmdBuilder(args...)
	}
}

func (database *Database) mysqldumpCommandBuilder(direction string, args ...string) []interface{} {
	if direction == "local" {
		return database.localMysqlCmdBuilder(args...)
	} else {
		return database.remoteMysqlCmdBuilder(args...)
	}
}

func (database *Database) mysqlTableList(connectionType string) []string {
	sqlStmt := "SHOW TABLES"

	cmd := shell.Cmd("echo", sqlStmt).Pipe(database.mysqlCommandBuilder(connectionType)...)
	output := cmd.Run().Stdout.String()

	outputString := strings.TrimSpace(string(output))
	ret := strings.Split(outputString, "\n")

	return ret
}
