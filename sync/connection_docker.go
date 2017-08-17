package sync

import (
	"strings"
	"fmt"
	"github.com/webdevops/go-shell"
)

var containerCache = map[string]string{}

func (connection *Connection) DockerCommandBuilder(cmd string, args ...string) []interface{} {
	dockerArgs := append(ConnectionDockerArguments, connection.DockerGetContainerId(), cmd)
	dockerArgs = append(dockerArgs, args...)

	if connection.GetType() == "ssh+docker" {
		return connection.SshCommandBuilder("docker", dockerArgs...)
	} else {
		return connection.LocalCommandBuilder("docker", dockerArgs...)
	}
}

func (connection *Connection) DockerGetContainerId() string {
	var container string

	cacheKey := fmt.Sprintf("%s:%s", connection.Hostname, connection.Docker)

	if val, ok := containerCache[cacheKey]; ok {
		// use cached
		container = val
	} else if strings.HasPrefix(connection.Docker, "compose:") {

		// copy connection because we need conn without docker usage (endless loop)
		connectionClone := *connection
		connectionClone.Docker = ""
		connectionClone.Type  = "auto"

		// docker-compose
		containerName := strings.TrimPrefix(connection.Docker, "compose:")

		cmd := shell.Cmd(connectionClone.CommandBuilder("docker-compose", "ps", "-q", containerName)...).Run()
		containerId := strings.TrimSpace(cmd.Stdout.String())

		if containerId == "" {
			panic(fmt.Sprintf("Container \"%s\" not found empty", container))
		}

		container = containerId
	} else {
		// normal docker
		container = connection.Docker
	}

	// cache value
	containerCache[cacheKey] = container

	return container
}
