package sync

import (
	"github.com/mohae/deepcopy"
	"github.com/webdevops/go-shell/commandbuilder"
)

type YamlCommandBuilderConnection struct {
	Type string
	Ssh *YamlCommandBuilderArgument
	Docker *YamlCommandBuilderArgument

	Environment *map[string]string
	Workdir string

	connection *commandbuilder.Connection
}

// Get (or create) connection instance
// will be cached one it's created
func (yconn *YamlCommandBuilderConnection) GetInstance() *commandbuilder.Connection {
	if yconn.connection == nil {
		conn := commandbuilder.Connection{}
		conn.Type = yconn.Type

		if yconn.Ssh != nil {
			conn.Ssh = yconn.Ssh.Argument
		}

		if yconn.Docker != nil {
			conn.Docker = yconn.Docker.Argument
		}

		if yconn.Environment != nil {
			conn.Environment.SetMap(*yconn.Environment)
		}

		if yconn.Workdir != "" {
			conn.Workdir = yconn.Workdir
		}

		yconn.connection = &conn
	}

	return yconn.connection
}

// Checks if connection is empty
func (yconn *YamlCommandBuilderConnection) IsEmpty() (status bool) {
	status = false
	if yconn.Type != ""    { return }
	if yconn.Ssh != nil    { return }
	if yconn.Docker != nil { return }
	if yconn.Environment != nil { return }

	return true
}

// Clone yaml connection (without shell connection instance)
func (yconn *YamlCommandBuilderConnection) Clone() (conn *YamlCommandBuilderConnection) {
	conn = deepcopy.Copy(yconn).(*YamlCommandBuilderConnection)
	conn.connection = nil
	return conn
}
