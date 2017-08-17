package sync

import "gopkg.in/yaml.v2"

func (server *Server) GetLocalPath() string {
	if server.Path == "" {
		Logger.FatalExit(1, "server.Path is empty")
	}

	return server.Path
}

func (server *Server) AsYaml() string {
	conf, _ := yaml.Marshal(server)
	return string(conf)
}
