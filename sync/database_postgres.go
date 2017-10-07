package sync

import (
	"fmt"
	"strings"
	"github.com/webdevops/go-shell"
)

type DatabasePostgres struct {
	Database
}


func (database *DatabasePostgres) init() {
	// LOCAL
	if database.Local.Connection.Docker != "" {
		queryConn := database.Local.Connection.Clone()
		queryConn.Type = "auto"
		queryConn.Docker = ""

		// docker auto hostname
		database.Local.Hostname = "127.0.0.1"

		if database.Local.User == "" || database.Local.Schema == "" {
			containerId := queryConn.DockerGetContainerId(database.Local.Connection.Docker)
			containerEnv := queryConn.DockerGetEnvironment(containerId)

			// try to guess user/password
			if database.Local.User == "" {
				// get superuser pass from env
				if pass, ok := containerEnv["POSTGRES_PASSWORD"]; ok {
					if user, ok := containerEnv["POSTGRES_USER"]; ok {
						fmt.Println("   -> local: using postgres superadmin account (from POSTGRES_USER and POSTGRES_PASSWORD)")
						database.Local.User = user
						database.Local.Password = pass
					} else {
						fmt.Println("   -> local: using postgres superadmin account (from POSTGRES_PASSWORD)")
						// only password available
						database.Local.User = "postgres"
						database.Local.Password = pass
					}
				}
			}

			// get schema from env
			if database.Local.Schema == "" {
				if val, ok := containerEnv["POSTGRES_DB"]; ok {
					fmt.Println("   -> local: using postgres schema (from POSTGRES_DB)")
					database.Local.Schema = val
				}
			}
		}
	}

	// Remote
	if database.Connection.Docker != "" {
		queryConn := database.Connection.Clone()
		queryConn.Type = "auto"
		queryConn.Docker = ""

		// docker auto hostname
		database.Hostname = "127.0.0.1"

		if database.User == "" || database.Schema == "" {
			containerId := queryConn.DockerGetContainerId(database.Connection.Docker)
			containerEnv := queryConn.DockerGetEnvironment(containerId)

			// try to guess user/password
			if database.User == "" {
				// get superuser pass from env
				if pass, ok := containerEnv["POSTGRES_PASSWORD"]; ok {
					if user, ok := containerEnv["POSTGRES_USER"]; ok {
						fmt.Println("   -> remote: using postgres superadmin account (from POSTGRES_USER and POSTGRES_PASSWORD)")
						database.User = user
						database.Password = pass
					} else {
						fmt.Println("   -> remote: using postgres superadmin account (from POSTGRES_PASSWORD)")
						// only password available
						database.User = "postgres"
						database.Password = pass
					}
				}
			}

			// get schema from env
			if database.Schema == "" {
				if val, ok := containerEnv["POSTGRES_DB"]; ok {
					fmt.Println("   -> remote: using postgres schema (from POSTGRES_DB)")
					database.Schema = val
				}
			}
		}
	}
}

func (database *DatabasePostgres) tableFilter(connectionType string) ([]string, []string) {
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
