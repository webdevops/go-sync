package sync

type server struct {
	Path string
	Connection connection
	Filesystem []filesystem
	Database []database
}

func (server *server) Sync() {
	defer func() {
		recover := recover()
		ShellErrorHandler(recover)
	}()

	server.SyncFilesystem()
	server.SyncDatabases()
}

func (server *server) SyncFilesystem() {
	for _, filesystem := range server.Filesystem {
		Logger.Main("Starting sync of %s", filesystem.String(server))
		filesystem.Sync(server)
	}
}

func (server *server) SyncDatabases() {
	for _, database := range server.Database {
		Logger.Main("Starting sync of %s", database.String())
		database.Sync(server)
	}
}

func (server *server) GetLocalPath() string {
	if server.Path == "" {
		Logger.FatalExit(1, "server.Path is empty")
	}

	return server.Path
}
