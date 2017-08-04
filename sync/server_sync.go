package sync

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
		Logger.Main("Starting sync of %s", filesystem.String(server, "sync"))
		filesystem.Sync(server)
	}
}

func (server *server) SyncDatabases() {
	for _, database := range server.Database {
		Logger.Main("Starting sync of %s", database.String("sync"))
		database.Sync(server)
	}
}
