package sync

import (
	"fmt"
	"strings"
)

type filesystem struct {
	Path string
	Local string
	Filter filter
}

func (filesystem *filesystem) localPath(server *server) string {
	if filesystem.Local != "" {
		return filesystem.Local
	} else {
		return server.GetLocalPath()
	}
}

func (filesystem *filesystem) String(server *server, direction string) string {
	var parts []string

	switch direction {
	case "sync":
		parts = append(parts, fmt.Sprintf("Path:%s", filesystem.Path))
		parts = append(parts, "->")
		parts = append(parts, fmt.Sprintf("Local:%s", filesystem.localPath(server)))
	case "deploy":
		parts = append(parts, fmt.Sprintf("Local:%s", filesystem.localPath(server)))
		parts = append(parts, "->")
		parts = append(parts, fmt.Sprintf("Path:%s", filesystem.Path))
	}

	return fmt.Sprintf("Filesystem[%s]", strings.Join(parts[:]," "))
}
