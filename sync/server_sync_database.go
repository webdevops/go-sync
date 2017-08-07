package sync

import (
	"github.com/webdevops/go-shell"
	"fmt"
)

func (database *Database) Sync(server *Server) {
	database.remoteConnection = server.Connection

	if database.Options.ClearDatabase {
		database.syncClearDatabase(server)
	}

	database.syncStructure(server)
	database.syncData(server)
}

// Sync database structure
func (database *Database) syncClearDatabase(server *Server) {

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
func (database *Database) syncStructure(server *Server) {
	Logger.Step("syncing database structure")

	// Sync structure only
	dumpCmd := database.remoteMysqldumpCmdBuilder([]string{"--no-data"}, false)
	restoreCmd := database.localMysqlCmdBuilder()

	cmd := shell.Cmd(dumpCmd...).Pipe("gunzip", "--stdout").Pipe(restoreCmd...)
	cmd.Run()
}

// Sync database data
func (database *Database) syncData(server *Server) {
	Logger.Step("syncing database data")

	// Sync data only
	dumpCmd := database.remoteMysqldumpCmdBuilder([]string{"--no-create-info"}, true)
	restoreCmd := database.localMysqlCmdBuilder()

	cmd := shell.Cmd(dumpCmd...).Pipe("gunzip", "--stdout").Pipe(restoreCmd...)
	cmd.Run()
}
