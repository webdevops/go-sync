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
	for _, filesystem := range server.Filesystem {
		Logger.Main("Starting deploy of %s", filesystem.String(server, "deploy"))
		filesystem.Deploy(server)
	}
}

func (server *Server) DeployDatabases() {
	for _, database := range server.Database {
		Logger.Main("Starting deploy of %s", database.String("deploy"))
		database.Deploy(server)
	}
}
