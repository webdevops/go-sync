package sync

import (
	"fmt"
	"strings"
	"github.com/webdevops/go-shell/commandbuilder"
)

func (filesystem *Filesystem) ApplyDefaults(server *Server) {
	// set default connection if not set
	if (commandbuilder.Connection{}) == filesystem.Connection {
		filesystem.Connection = server.Connection
	}

	// set default path
	if filesystem.Local == "" {
		filesystem.Local = server.GetLocalPath()
	}
}

func (filesystem *Filesystem) localPath() string {
	return filesystem.Local
}

func (filesystem *Filesystem) String(direction string) string {
	var parts []string

	switch direction {
	case "sync":
		parts = append(parts, fmt.Sprintf("Path:%s", filesystem.Path))
		parts = append(parts, "->")
		parts = append(parts, fmt.Sprintf("Local:%s", filesystem.localPath()))
	case "deploy":
		parts = append(parts, fmt.Sprintf("Local:%s", filesystem.localPath()))
		parts = append(parts, "->")
		parts = append(parts, fmt.Sprintf("Path:%s", filesystem.Path))
	}

	return fmt.Sprintf("Filesystem[%s]", strings.Join(parts[:]," "))
}
