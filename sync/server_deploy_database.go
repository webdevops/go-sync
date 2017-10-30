package sync

import (
	"github.com/webdevops/go-shell"
	"fmt"
)

func (database *Database) Deploy() {
	switch database.GetType() {
	case "mysql":
		mysql := database.GetMysql()
		if mysql.Options.ClearDatabase {
			mysql.deployClearDatabase()
		}

		mysql.deployStructure()
		mysql.deployData()

	case "postgres":
		postgres := database.GetPostgres()
		fmt.Sprintf(postgres.String("deploy"))
	}
}

// Deploy database structure
func (database *DatabaseMysql) deployClearDatabase() {

	// don't use database which we're trying to drop, instead use "mysql"
	db := database.Db
	database.Db = ""

	Logger.Step("dropping remote database \"%s\"", db)
	dropStmt := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", db)
	dropCmd := shell.Cmd("echo", shell.Quote(dropStmt)).Pipe(database.remoteMysqlCmdBuilder()...)
	dropCmd.Run()

	Logger.Step("creating remote database \"%s\"", db)
	createStmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", db)
	createCmd := shell.Cmd("echo", shell.Quote(createStmt)).Pipe(database.remoteMysqlCmdBuilder()...)
	createCmd.Run()

	database.Db = db
}

// Deploy database structure
func (database *DatabaseMysql) deployStructure() {
	Logger.Step("deploy database structure")

	// Deploy structure only
	dumpCmd := database.localMysqldumpCmdBuilder([]string{"--no-data"}, false)
	restoreCmd := database.remoteMysqlCmdBuilderUncompress()

	cmd := shell.Cmd(dumpCmd...).Pipe("gzip", "--stdout").Pipe(restoreCmd...)
	cmd.Run()
}

// Deploy database data
func (database *DatabaseMysql) deployData() {
	Logger.Step("deploy database data")

	// Deploy data only
	dumpCmd := database.localMysqldumpCmdBuilder([]string{"--no-create-info"}, true)
	restoreCmd := database.remoteMysqlCmdBuilderUncompress()

	cmd := shell.Cmd(dumpCmd...).Pipe("gzip", "--stdout").Pipe(restoreCmd...)
	cmd.Run()
}
