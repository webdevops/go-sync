package sync

func (server *Server) RunExec(when string) {
	defer func() {
		recover := recover()
		ShellErrorHandler(recover)
	}()

	execList := server.GetExecByWhen(when)

	if len(execList) >= 1 {
		Logger.Main("Starting exec mode \"%s\"", when)

		for _, exec := range execList {
			Logger.Step("executing >> %s", exec.String(server))
			exec.Execute(server)
		}
	}
}

func (server *Server) GetExecByWhen(when string) []Execution {
	var execList []Execution

	for _, val := range server.Exec {
		execList = append(execList, val)
	}

	return execList
}
