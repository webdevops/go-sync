package sync

import (
	"fmt"
	"strings"
)

type database struct {
	Type string
	Schema string
	Hostname string
	Port string
	User string
	Password string
	Filter filter
	Local struct {
		Schema string
		Hostname string
		Port string
		User string
		Password string
		Connection connection
	}
	Options struct {
		ClearDatabase bool `yaml:"clear-database"`
	}

	// local cache
	cacheRemoteTableList []string
	cacheLocalTableList []string

	remoteConnection connection
}

func (database *database) mysqlTableFilter(connection *connection, direction string) ([]string, []string) {
	var exclude []string
	var include []string

	var tableList []string

	if (direction == "sync") {
		if len(database.cacheLocalTableList) == 0 {
			Logger.Step("get list of mysql tables for table filter")
			database.cacheLocalTableList = database.localMysqlTableList()
		}

		tableList = database.cacheRemoteTableList

	} else {
		if len(database.cacheRemoteTableList) == 0 {
			Logger.Step("get list of mysql tables for table filter")
			database.cacheRemoteTableList = database.remoteMysqlTableList()
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

func (database *database) String(direction string) string {
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
