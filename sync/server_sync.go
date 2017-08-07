package sync

func (server *Server) Sync() {
	defer func() {
		recover := recover()
		ShellErrorHandler(recover)
	}()

	server.SyncFilesystem()
	server.SyncDatabases()
}

func (server *Server) SyncFilesystem() {
	for _, filesystem := range server.Filesystem {
		Logger.Main("Starting sync of %s", filesystem.String(server, "sync"))
		filesystem.Sync(server)
	}
}

func (server *Server) SyncDatabases() {
	for _, database := range server.Database {
		Logger.Main("Starting sync of %s", database.String("sync"))
		database.Sync(server)
	}
}
