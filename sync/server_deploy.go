package sync

func (server *server) Deploy() {
	defer func() {
		recover := recover()
		ShellErrorHandler(recover)
	}()

	server.DeployFilesystem()
	server.DeployDatabases()
}

func (server *server) DeployFilesystem() {
	for _, filesystem := range server.Filesystem {
		Logger.Main("Starting deploy of %s", filesystem.String(server, "deploy"))
		filesystem.Deploy(server)
	}
}

func (server *server) DeployDatabases() {
	for _, database := range server.Database {
		Logger.Main("Starting deploy of %s", database.String("deploy"))
		database.Deploy(server)
	}
}
