package sync

import (
	"fmt"
	"strings"
	"github.com/mohae/deepcopy"
)

func (database *Database) ApplyDefaults(server *Server) {
	// set default connection if not set
	if database.Connection.IsEmpty() {
		database.Connection = deepcopy.Copy(server.Connection).(YamlCommandBuilderConnection)
	}
}

func (database *Database) GetType() string {
	var dbtype string

	switch database.Type {
	case "mysql":
		dbtype = "mysql"
	case "postgresql":
		fallthrough
	case "postgres":
		dbtype = "postgres"
	default:
		panic(fmt.Sprintf("Database type %s is not valid or supported", database.Type))
	}

	return dbtype
}

func (database *Database) GetMysql() DatabaseMysql {
	mysql := DatabaseMysql{*database}
	mysql.init()
	return mysql
}

func (database *Database) GetPostgres() DatabasePostgres {
	postgres := DatabasePostgres{*database}
	postgres.init()
	return postgres
}

func (database *Database) String(direction string) string {
	var parts, remote, local []string
	
	connRemote := database.Connection.GetInstance()
	connLocal := database.Local.Connection.GetInstance()

	// general
	parts = append(parts, fmt.Sprintf("Type:%s", database.Type))

	//-------------------------------------------
	// remote
	remote = append(remote, fmt.Sprintf("Schema:%s", database.Schema))
	remote = append(remote, fmt.Sprintf("Connection:%s", connRemote.GetType()))

	if connRemote.IsSsh() {
		if connRemote.SshConnectionHostnameString() != "" {
			remote = append(remote, fmt.Sprintf("SSH:%s", connRemote.SshConnectionHostnameString()))
		}
	}

	if connRemote.IsDocker() {
		remote = append(remote, fmt.Sprintf("Docker:%s", connRemote.Docker.Hostname))
	} else if database.Hostname != "" {
		hostname := database.Hostname

		if database.Port != "" {
			hostname += ":"+database.Port
		}
		remote = append(remote, fmt.Sprintf("Host:%s", hostname))
	}

	if database.User != "" {
		remote = append(remote, fmt.Sprintf("User:%s", database.User))
	}

	if database.Password != "" {
		remote = append(remote, fmt.Sprintf("Passwd:%s", "*****"))
	}

	//-------------------------------------------
	// local
	local = append(local, fmt.Sprintf("Schema:%s", database.Local.Schema))
	local = append(local, fmt.Sprintf("Connection:%s", connLocal.GetType()))

	if connLocal.IsSsh() {
		if connLocal.SshConnectionHostnameString() != "" {
			local = append(local, fmt.Sprintf("SSH:%s", connLocal.SshConnectionHostnameString()))
		}
	}

	if connLocal.IsDocker() {
		local = append(local, fmt.Sprintf("Docker:%s", connLocal.Docker.Hostname))
	} else if database.Local.Hostname != "" {
		hostname := database.Local.Hostname

		if database.Local.Port != "" {
			hostname += ":"+database.Local.Port
		}
		local = append(local, fmt.Sprintf("Host:%s", hostname))
	}

	// build parts
	switch direction {
	case "sync":
		parts = append(parts, remote...)
		parts = append(parts, "->")
		parts = append(parts, local...)
	case "deploy":
		parts = append(parts, local...)
		parts = append(parts, "->")
		parts = append(parts, remote...)
	}

	return fmt.Sprintf("Database[%s]", strings.Join(parts[:]," "))
}
