package sync

func (server *Server) Sync() {
	defer func() {
		recover := recover()
		ShellErrorHandler(recover)
	}()

	server.Init()

	if server.runConfiguration.Exec {
		server.RunExec("startup")
	}

	if server.runConfiguration.Filesystem {
		server.SyncFilesystem()
	}

	if server.runConfiguration.Database {
		server.SyncDatabases()
	}

	if server.runConfiguration.Exec {
		server.RunExec("finish")
	}

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
