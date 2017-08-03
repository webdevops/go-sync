package sync

import (
	"fmt"
	"strings"
	"github.com/webdevops/go-shell"
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
	remoteTableList []string

	remoteConnection connection
}

func (database *database) localCommandInterface(command string, args ...string) []interface{} {
	var ret []interface{}

	if database.Local.Connection.Type == "" {
		database.Local.Connection.Type = "local"

		// autodetection
		if database.Local.Connection.Docker != "" {
			database.Local.Connection.Type = "docker"
		}

		if database.Local.Connection.Hostname != "" {
			database.Local.Connection.Type = "ssh"
		}
	}

	switch database.Local.Connection.Type {
	case "local":
		ret = ShellCommandInterfaceBuilder(command, args...)
	case "ssh":
		ret = database.Local.Connection.RemoteCommandBuilder(command, args...)
	}

	return ret
}

func (database *database) remoteSshDump(additionalArgs []string, useFilter bool) []interface{} {
	var args []string

	if database.User != "" {
		args = append(args, "-u" + database.User)
	}

	if database.Password != "" {
		args = append(args, "-p" + database.Password)
	}

	if database.Hostname != "" {
		args = append(args, "-h" + database.Hostname)
	}

	if database.Port != "" {
		args = append(args, "-P" + database.Port)
	}

	if len(args) > 0 {
		args = append(args, additionalArgs...)
	}

	// exclude
	excludeArgs, includeArgs := database.sshFilterDump(&database.remoteConnection);
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// schema
	args = append(args, database.Schema)

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	return database.remoteConnection.SshCompressedCommandBuilder("mysqldump", args...)
}

func (database *database) sshFilterDump(connection *connection) ([]string, []string) {
	var exclude []string
	var include []string

	if len(database.remoteTableList) == 0 {
		Logger.Step("get list of mysql tables for table filter")
		database.remoteTableList = database.mysqlTableList()
	}

	// calc excludes
	excludeTableList := database.Filter.CalcExcludes(database.remoteTableList)
	for _, table := range excludeTableList {
		exclude  = append(exclude, fmt.Sprintf("--ignore-table=%s.%s", database.Schema, table))
	}

	// calc includes
	includeTableList := database.Filter.CalcIncludes(database.remoteTableList)
	for _, table := range includeTableList {
		include  = append(include, table)
	}

	return exclude, include
}

func (database *database) remoteMysqlCmdBuilder(args ...string) []interface{} {
	args = append(args, "-BN")

	if database.User != "" {
		args = append(args, "-u" + database.User)
	}

	if database.Password != "" {
		args = append(args, "-p" + database.Password)
	}

	if database.Hostname != "" {
		args = append(args, "-h" + database.Hostname)
	}

	if database.Port != "" {
		args = append(args, "-P" + database.Port)
	}

	args = append(args, database.Schema)

	return database.remoteConnection.RemoteCommandBuilder("mysql", args...)
}

func (database *database) mysqlCmdBuilderLocal(args ...string) []interface{} {
	args = append(args, "-BN")

	if database.Local.User != "" {
		args = append(args, "-u" + database.Local.User)
	}

	if database.Local.Password != "" {
		args = append(args, "-p" + database.Local.Password)
	}

	if database.Local.Hostname != "" {
		args = append(args, "-h" + database.Local.Hostname)
	}

	if database.Local.Port != "" {
		args = append(args, "-P" + database.Local.Port)
	}

	args = append(args, database.Local.Schema)

	return database.Local.Connection.RemoteCommandBuilder("mysql", args...)
}

func (database *database) mysqlTableList() []string {
	sqlStmt := "SHOW TABLES"
	cmd := shell.Cmd("echo", sqlStmt).Pipe(database.remoteMysqlCmdBuilder()...)
	output := cmd.Run().Stdout.String()

	outputString := strings.TrimSpace(string(output))
	ret := strings.Split(outputString, "\n")

	return ret
}

func (database *database) String() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Schema:%s", database.Schema))

	if database.Hostname != "" {
		parts = append(parts, fmt.Sprintf("Host:%s", database.Hostname))
	}

	if database.User != "" {
		parts = append(parts, fmt.Sprintf("User:%s", database.User))
	}

	if database.Password != "" {
		parts = append(parts, fmt.Sprintf("Passwd:%s", "*****"))
	}

	parts = append(parts, "->")

	parts = append(parts, fmt.Sprintf("Schema:%s", database.Local.Schema))


	return fmt.Sprintf("Database[%s]", strings.Join(parts[:]," "))
}
