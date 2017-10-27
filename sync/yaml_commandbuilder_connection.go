package sync

import (
	"github.com/webdevops/go-shell/commandbuilder"
)

type YamlCommandBuilderConnection struct {
	Type string
	Ssh YamlCommandBuilderArgument
	Docker YamlCommandBuilderArgument
}

func (yconn *YamlCommandBuilderConnection) GetInstance() *commandbuilder.Connection {
	conn := commandbuilder.Connection{}
	conn.Type = yconn.Type

	conn.Ssh = yconn.Ssh.Argument
	conn.Docker = yconn.Docker.Argument

	return &conn
}

func (yconn *YamlCommandBuilderConnection) IsEmpty() bool {
	return yconn.Type == "" && yconn.Docker.IsEmpty() && yconn.Ssh.IsEmpty()
}
