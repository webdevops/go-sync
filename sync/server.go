package sync

import (
	"gopkg.in/yaml.v2"
)

func (server *Server) Init() {
	if server.runConfiguration == nil {
		server.runConfiguration = &RunConfiguration{
			Database: true,
			Filesystem: true,
		}
	}
}

func (server *Server) GetLocalPath() string {
	if server.Path == "" {
		Logger.FatalExit(1, "server.Path is empty")
	}

	return server.Path
}

func (server *Server) SetRunConfiguration(conf RunConfiguration) {
	server.runConfiguration = &conf
}

func (server *Server) AsYaml() string {
	conf, _ := yaml.Marshal(server)
	return string(conf)
}
