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
		if filesystem.Options.GenerateStubs {
			Logger.Main("Starting stub generator for %s", filesystem.String(server, "sync"))
			filesystem.SyncStubs(server)
		} else {
			Logger.Main("Starting sync of %s", filesystem.String(server, "sync"))
			filesystem.Sync(server)
		}
	}
}

func (server *Server) SyncDatabases() {
	for _, database := range server.Database {
		Logger.Main("Starting sync of %s", database.String("sync"))
		database.Sync(server)
	}
}
