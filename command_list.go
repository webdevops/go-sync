package main

type ListCommand struct {
	AbstractCommand
}

// List all possible server configurations (for sync and deploy)
func (command *ListCommand) Execute(args []string) error {
	config := command.GetConfig()
	config.ShowConfiguration()
	return nil
}
