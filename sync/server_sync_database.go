package sync

import (
	"github.com/webdevops/go-shell"
	"fmt"
	"io/ioutil"
	"os"
)

func (database *Database) Sync() {
	switch database.GetType() {
	case "mysql":
		mysql := database.GetMysql()
		if mysql.Options.ClearDatabase {
			mysql.syncClearDatabase()
		}

		mysql.syncStructure()
		mysql.syncData()

	case "postgres":
		postgres := database.GetPostgres()

		if postgres.Options.ClearDatabase {
			postgres.syncClearDatabase()
		}

		postgres.syncStructure()
		postgres.syncData()
	}
}

//#############################################################################
// Postgres
//#############################################################################

// Sync database structure
func (database *DatabasePostgres) syncClearDatabase() {

	// don't use database which we're trying to drop, instead use "mysql"
	db := database.Local.Db
	database.Local.Db = "postgres"

	Logger.Step("dropping local database \"%s\"", db)
	dropStmt := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", db)
	dropCmd := shell.Cmd("echo", shell.Quote(dropStmt)).Pipe(database.localPsqlCmdBuilder()...)
	dropCmd.Run()

	Logger.Step("creating local database \"%s\"", db)
	createStmt := fmt.Sprintf("CREATE DATABASE `%s`", db)
	createCmd := shell.Cmd("echo", shell.Quote(createStmt)).Pipe(database.localPsqlCmdBuilder()...)
	createCmd.Run()

	database.Local.Db = db
}

// Sync database structure
func (database *DatabasePostgres) syncStructure() {
	Logger.Step("syncing database structure")

	tmpfile, err := ioutil.TempFile("", "dump")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())

	// Sync structure only
	dumpCmd := database.remotePgdumpCmdBuilder([]string{"--schema-only"}, false)
	shell.Cmd(dumpCmd...).Pipe("cat", ">", tmpfile.Name()).Run()

	// Restore structure only
	restoreCmd := database.localPsqlCmdBuilder()
	shell.Cmd("cat", tmpfile.Name()).Pipe("gunzip", "--stdout").Pipe(restoreCmd...).Run()
}


// Sync database data
func (database *DatabasePostgres) syncData() {
	Logger.Step("syncing database data")

	tmpfile, err := ioutil.TempFile("", "dump")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())

	// Sync structure only
	dumpCmd := database.remotePgdumpCmdBuilder([]string{"--data-only"}, true)
	shell.Cmd(dumpCmd...).Pipe("cat", ">", tmpfile.Name()).Run()

	// Restore structure only
	restoreCmd := database.localPsqlCmdBuilder()
	shell.Cmd("cat", tmpfile.Name()).Pipe("gunzip", "--stdout").Pipe(restoreCmd...).Run()
}

//#############################################################################
// MySQL
//#############################################################################

// Sync database structure
func (database *DatabaseMysql) syncClearDatabase() {

	// don't use database which we're trying to drop, instead use "mysql"
	db := database.Local.Db
	database.Local.Db = ""

	Logger.Step("dropping local database \"%s\"", db)
	dropStmt := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", db)
	dropCmd := shell.Cmd("echo", shell.Quote(dropStmt)).Pipe(database.localMysqlCmdBuilder()...)
	dropCmd.Run()

	Logger.Step("creating local database \"%s\"", db)
	createStmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", db)
	createCmd := shell.Cmd("echo", shell.Quote(createStmt)).Pipe(database.localMysqlCmdBuilder()...)
	createCmd.Run()

	database.Local.Db = db
}

// Sync database structure
func (database *DatabaseMysql) syncStructure() {
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
func (database *DatabaseMysql) syncData() {
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
