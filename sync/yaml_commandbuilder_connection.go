package sync

import (
	"github.com/webdevops/go-shell/commandbuilder"
)

type YamlCommandBuilderConnection struct {
	Type string
	Ssh YamlCommandBuilderArgument
	Docker YamlCommandBuilderArgument

	connection *commandbuilder.Connection
}

func (yconn *YamlCommandBuilderConnection) GetInstance() *commandbuilder.Connection {
	if yconn.connection == nil {
		conn := commandbuilder.Connection{}
		conn.Type = yconn.Type
		conn.Ssh = yconn.Ssh.Argument
		conn.Docker = yconn.Docker.Argument

		yconn.connection = &conn
	}

	return yconn.connection
}

func (yconn *YamlCommandBuilderConnection) IsEmpty() bool {
	return yconn.Type == "" && yconn.Docker.IsEmpty() && yconn.Ssh.IsEmpty()
}
