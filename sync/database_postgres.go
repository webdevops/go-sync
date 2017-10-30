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
	connLocal := database.Local.Connection.GetInstance()
	connRemote := database.Connection.GetInstance()

	// LOCAL
	if connLocal.IsDocker() {
		if database.Local.User == "" || database.Local.Db == "" {
			containerEnv := connLocal.DockerGetEnvironment()

			// try to guess user/password
			if database.Local.User == "" {
				// get superuser pass from env
				if pass, ok := containerEnv["POSTGRES_PASSWORD"]; ok {
					if user, ok := containerEnv["POSTGRES_USER"]; ok {
						fmt.Println(fmt.Sprintf("   -> local: using postgres superadmin account \"%s\" (from env:POSTGRES_USER and env:POSTGRES_PASSWORD)", user))
						database.Local.User = user
						database.Local.Password = pass
					} else {
						fmt.Println("   -> local: using postgres superadmin account \"postgres\" (from env:POSTGRES_PASSWORD)")
						// only password available
						database.Local.User = "postgres"
						database.Local.Password = pass
					}
				}
			}

			// get database from env
			if database.Local.Db == "" {
				if db, ok := containerEnv["POSTGRES_DB"]; ok {
					fmt.Println(fmt.Sprintf("   -> remote: using postgres database \"%s\" (from env:POSTGRES_DB)", db))
					database.Local.Db = db
				}
			}
		}
	}

	// Remote
	if connRemote.IsDocker() {
		if database.User == "" || database.Db == "" {
			containerEnv := connRemote.DockerGetEnvironment()

			// try to guess user/password
			if database.User == "" {
				// get superuser pass from env
				if pass, ok := containerEnv["POSTGRES_PASSWORD"]; ok {
					if user, ok := containerEnv["POSTGRES_USER"]; ok {
						fmt.Println(fmt.Sprintf("   -> remote: using postgres superadmin account \"%s\" (from env:POSTGRES_USER and env:POSTGRES_PASSWORD)", user))
						database.User = user
						database.Password = pass
					} else {
						fmt.Println("   -> remote: using postgres superadmin account \"postgres\" (from env:POSTGRES_PASSWORD)")
						// only password available
						database.User = "postgres"
						database.Password = pass
					}
				}
			}

			// get database from env
			if database.Db == "" {
				if db, ok := containerEnv["POSTGRES_DB"]; ok {
					fmt.Println(fmt.Sprintf("   -> remote: using postgres database \"%s\" (from env:POSTGRES_DB)", db))
					database.Db = db
				}
			}
		}
	}
}

func (database *DatabasePostgres) tableFilter(connectionType string) (exclude []string, include []string) {
	var tableList []string

	if connectionType == "local" {
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
	
	return
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
	sqlStmt = fmt.Sprintf(sqlStmt, database.quote(database.Db))

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
