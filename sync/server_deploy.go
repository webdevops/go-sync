package sync

func (server *Server) Deploy() {
	defer func() {
		recover := recover()
		ShellErrorHandler(recover)
	}()

	server.RunExec("startup")
	server.DeployFilesystem()
	server.DeployDatabases()
	server.RunExec("finish")

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
