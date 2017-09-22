package sync

import "fmt"

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
	switch when {
	case "startup":
		return server.ExecStartup
	case "finish":
		return server.ExecFinish
	default:
		panic(fmt.Sprintf("execution list %s is not valid", when))
	}
}
