package sync

func (server *Server) Sync() {
	defer func() {
		recover := recover()
		ShellErrorHandler(recover)
	}()

	server.RunExec("startup")
	server.SyncFilesystem()
	server.SyncDatabases()
	server.RunExec("finish")

	waitGroup.Wait()
}

func (server *Server) SyncFilesystem() {
	for _, filesystem := range server.Filesystem {
		filesystem.ApplyDefaults(server)

		if filesystem.Options.GenerateStubs {
			Logger.Main("Starting stub generator for %s", filesystem.String( "sync"))
			filesystem.SyncStubs()
		} else {
			Logger.Main("Starting sync of %s", filesystem.String("sync"))
			filesystem.Sync()
		}
	}
}

func (server *Server) SyncDatabases() {
	for _, database := range server.Database {
		database.ApplyDefaults(server)
		Logger.Main("Starting sync of %s", database.String("sync"))
		database.Sync()
	}
}
