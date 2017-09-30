package sync

import (
	"github.com/webdevops/go-shell"
	"fmt"
	"io/ioutil"
	"os"
)

func (database *Database) Sync() {
	if database.Options.ClearDatabase {
		database.syncClearDatabase()
	}

	database.syncStructure()
	database.syncData()
}

// Sync database structure
func (database *Database) syncClearDatabase() {

	// don't use database which we're trying to drop, instead use "mysql"
	schema := database.Local.Schema
	database.Local.Schema = ""

	Logger.Step("dropping local database \"%s\"", schema)
	dropStmt := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", schema)
	dropCmd := shell.Cmd("echo", shell.Quote(dropStmt)).Pipe(database.localMysqlCmdBuilder()...)
	dropCmd.Run()

	Logger.Step("creating local database \"%s\"", schema)
	createStmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", schema)
	createCmd := shell.Cmd("echo", shell.Quote(createStmt)).Pipe(database.localMysqlCmdBuilder()...)
	createCmd.Run()

	database.Local.Schema = schema
}

// Sync database structure
func (database *Database) syncStructure() {
	Logger.Step("syncing database structure")

	tmpfile, err := ioutil.TempFile("", "dump")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())

	// Sync structure only
	dumpCmd := database.remoteMysqldumpCmdBuilder([]string{"--no-data"}, false)
	shell.Cmd(dumpCmd...).Pipe("cat", ">", tmpfile.Name()).Run()

	// Restore structure only
	restoreCmd := database.localMysqlCmdBuilder()
	shell.Cmd("cat", tmpfile.Name()).Pipe("gunzip", "--stdout").Pipe(restoreCmd...).Run()
}

// Sync database data
func (database *Database) syncData() {
	Logger.Step("syncing database data")

	tmpfile, err := ioutil.TempFile("", "dump")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())

	// Sync data only
	dumpCmd := database.remoteMysqldumpCmdBuilder([]string{"--no-create-info"}, false)
	shell.Cmd(dumpCmd...).Pipe("cat", ">", tmpfile.Name()).Run()

	// Restore data only
	restoreCmd := database.localMysqlCmdBuilder()
	shell.Cmd("cat", tmpfile.Name()).Pipe("gunzip", "--stdout").Pipe(restoreCmd...).Run()
}
