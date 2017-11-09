package sync

func (server *Server) Deploy() {
	defer func() {
		recover := recover()
		ShellErrorHandler(recover)
	}()

	server.Init()

	if server.runConfiguration.Exec {
		server.RunExec("startup")
	}

	if server.runConfiguration.Filesystem {
		server.DeployFilesystem()
	}

	if server.runConfiguration.Database {
		server.DeployDatabases()
	}

	if server.runConfiguration.Exec {
		server.RunExec("finish")
	}

	waitGroup.Wait()
}

func (server *Server) DeployFilesystem() {
	// check for generate-stubs option (not allowed)
	for _, filesystem := range server.Filesystem {
		if filesystem.Options.GenerateStubs {
			Logger.FatalExit(2, "Generate Stubs is not allowed for deployment")
		}
	}

	for _, filesystem := range server.Filesystem {
		filesystem.ApplyDefaults(server)

		Logger.Main("Starting deploy of %s", filesystem.String( "deploy"))
		filesystem.Deploy()
	}
}

func (server *Server) DeployDatabases() {
	for _, database := range server.Database {
		database.ApplyDefaults(server)
		Logger.Main("Starting deploy of %s", database.String("deploy"))
		database.Deploy()
	}
}
