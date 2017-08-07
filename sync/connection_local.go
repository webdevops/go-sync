package sync

func (connection *Connection) LocalCommandBuilder(cmd string, args ...string) []interface{} {
	return ShellCommandInterfaceBuilder(cmd, args...)
}
