package sync

func (server *Server) GetLocalPath() string {
	if server.Path == "" {
		Logger.FatalExit(1, "server.Path is empty")
	}

	return server.Path
}
