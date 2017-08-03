package sync

import (
	"github.com/webdevops/go-shell"
	"fmt"
)

func (database *database) Sync(server *server) {
	database.remoteConnection = server.Connection

	if database.Options.ClearDatabase {
		database.syncClearDatabase(server)
	}

	database.syncStructure(server)
	database.syncData(server)
}

// Sync database structure
func (database *database) syncClearDatabase(server *server) {

	// don't use database which we're trying to drop, instead use "mysql"
	schema := database.Local.Schema
	database.Local.Schema = "mysql"


	Logger.Step("dropping local database \"%s\"", schema)
	dropStmt := fmt.Sprintf("DROP DATABASE IF EXISTS %s", schema)
	dropCmd := shell.Cmd("echo", dropStmt).Pipe(database.mysqlCmdBuilderLocal()...)
	dropCmd.Run()

	Logger.Step("creating local database \"%s\"", schema)
	createStmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS  %s", schema)
	createCmd := shell.Cmd("echo", createStmt).Pipe(database.mysqlCmdBuilderLocal()...)
	createCmd.Run()

	database.Local.Schema = schema
}

// Sync database structure
func (database *database) syncStructure(server *server) {
	Logger.Step("syncing database structure")

	// Sync structure only
	dumpCmd := database.remoteSshDump([]string{"--no-data"}, false)
	restoreCmd := database.mysqlCmdBuilderLocal()

	cmd := shell.Cmd(dumpCmd...).Pipe("gunzip", "--stdout").Pipe(restoreCmd...)
	cmd.Run()
}

// Sync database data
func (database *database) syncData(server *server) {
	Logger.Step("syncing database data")

	// Sync data only
	dumpCmd := database.remoteSshDump([]string{"--no-create-info"}, true)
	restoreCmd := database.mysqlCmdBuilderLocal()

	cmd := shell.Cmd(dumpCmd...).Pipe("gunzip", "--stdout").Pipe(restoreCmd...)
	cmd.Run()
}
