package sync

import (
	"github.com/webdevops/go-shell"
	"fmt"
)

func (database *Database) Deploy() {
	if database.Options.ClearDatabase {
		database.deployClearDatabase()
	}

	database.deployStructure()
	database.deployData()
}

// Deploy database structure
func (database *Database) deployClearDatabase() {

	// don't use database which we're trying to drop, instead use "mysql"
	schema := database.Schema
	database.Schema = ""

	Logger.Step("dropping remote database \"%s\"", schema)
	dropStmt := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", schema)
	dropCmd := shell.Cmd("echo", shell.Quote(dropStmt)).Pipe(database.remoteMysqlCmdBuilder()...)
	dropCmd.Run()

	Logger.Step("creating remote database \"%s\"", schema)
	createStmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", schema)
	createCmd := shell.Cmd("echo", shell.Quote(createStmt)).Pipe(database.remoteMysqlCmdBuilder()...)
	createCmd.Run()

	database.Schema = schema
}

// Deploy database structure
func (database *Database) deployStructure() {
	Logger.Step("deploy database structure")

	// Deploy structure only
	dumpCmd := database.localMysqldumpCmdBuilder([]string{"--no-data"}, false)
	restoreCmd := database.remoteMysqlCmdBuilderUncompress()

	cmd := shell.Cmd(dumpCmd...).Pipe("gzip", "--stdout").Pipe(restoreCmd...)
	cmd.Run()
}

// Deploy database data
func (database *Database) deployData() {
	Logger.Step("deploy database data")

	// Deploy data only
	dumpCmd := database.localMysqldumpCmdBuilder([]string{"--no-create-info"}, true)
	restoreCmd := database.remoteMysqlCmdBuilderUncompress()

	cmd := shell.Cmd(dumpCmd...).Pipe("gzip", "--stdout").Pipe(restoreCmd...)
	cmd.Run()
}
