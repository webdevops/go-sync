package sync

type server struct {
	Path string
	Connection connection
	Filesystem []filesystem
	Database []database
}

func (server *server) GetLocalPath() string {
	if server.Path == "" {
		Logger.FatalExit(1, "server.Path is empty")
	}

	return server.Path
}
