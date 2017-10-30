package sync

import (
	"fmt"
	"strings"
	"github.com/webdevops/go-shell"
)

type DatabaseMysql struct {
	Database
}

func (database *DatabaseMysql) init() {
	connLocal := database.Local.Connection.GetInstance()
	connRemote := database.Connection.GetInstance()
	
	// LOCAL
	if connLocal.IsDocker() {
		if database.Local.User == "" || database.Local.Db == "" {
			containerEnv := connLocal.DockerGetEnvironment()

			// try to guess user/password
			if database.Local.User == "" {
				if val, ok := containerEnv["MYSQL_ROOT_PASSWORD"]; ok {
					// get root pass from env
					if database.Local.User == "" && database.Local.Password == "" {
						fmt.Println("   -> local: using mysql root account (from env:MYSQL_ROOT_PASSWORD)")
						database.Local.User = "root"
						database.Local.Password = val
					}
				} else if val, ok := containerEnv["MYSQL_ALLOW_EMPTY_PASSWORD"]; ok {
					// get root without password from env
					if val == "yes" && database.Local.User == "" {
						fmt.Println("   -> local: using mysql root account (from env:MYSQL_ALLOW_EMPTY_PASSWORD)")
						database.Local.User = "root"
						database.Local.Password = ""
					}
				} else if user, ok := containerEnv["MYSQL_USER"]; ok {
					if pass, ok := containerEnv["MYSQL_PASSWORD"]; ok {
						if database.Local.User == "" && database.Local.Password == "" {
							fmt.Println(fmt.Sprintf("   -> local: using mysql user account \"%s\" (from env:MYSQL_USER and env:MYSQL_PASSWORD)", user))
							database.Local.User = user
							database.Local.Password = pass
						}
					}
				}
			}

			// get database from env
			if database.Local.Db == "" {
				if db, ok := containerEnv["MYSQL_DATABASE"]; ok {
					fmt.Println(fmt.Sprintf("   -> local: using mysql database \"%s\" (from env:MYSQL_DATABASE)", db))
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
				if val, ok := containerEnv["MYSQL_ROOT_PASSWORD"]; ok {
					// get root pass from env
					if database.User == "" && database.Password == "" {
						fmt.Println("   -> remote: using mysql root account (from env:MYSQL_ROOT_PASSWORD)")
						database.User = "root"
						database.Password = val
					}
				} else if val, ok := containerEnv["MYSQL_ALLOW_EMPTY_PASSWORD"]; ok {
					// get root without password from env
					if val == "yes" && database.User == "" {
						fmt.Println("   -> remote: using mysql root account (from env:MYSQL_ALLOW_EMPTY_PASSWORD)")
						database.User = "root"
						database.Password = ""
					}
				} else if user, ok := containerEnv["MYSQL_USER"]; ok {
					if pass, ok := containerEnv["MYSQL_PASSWORD"]; ok {
						if database.User == "" && database.Password == "" {
							fmt.Println(fmt.Sprintf("   -> remote: using mysql user account \"%s\" (from env:MYSQL_USER and env:MYSQL_PASSWORD)", user))
							database.User = user
							database.Password = pass
						}
					}
				}
			}

			// get database from env
			if database.Db == "" {
				if db, ok := containerEnv["MYSQL_DATABASE"]; ok {
					fmt.Println(fmt.Sprintf("   -> remote: using mysql database \"%s\" (from env:MYSQL_DATABASE)", db))
					database.Db = db
				}
			}
		}
	}
}

func (database *DatabaseMysql) tableFilter(connectionType string) (exclude []string, include []string) {
	var tableList []string

	if connectionType == "local" {
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
		exclude  = append(exclude, shell.Quote(fmt.Sprintf("--ignore-table=%s.%s", database.Db, table)))
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
