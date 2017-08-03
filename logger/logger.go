package logger

import (
	"log"
	"os"
	"fmt"
	"strings"
)

const (
	LogPrefix = ""
)

type SyncLogger struct {
	*log.Logger
}

var (
	Logger *SyncLogger
	Verbose bool
	CommandName string
)

func GetInstance(commandName string, flags int) *SyncLogger {
	CommandName = commandName

	if Logger == nil {
		Logger = &SyncLogger{log.New(os.Stdout, LogPrefix, flags)}
	}
	return Logger
}

func (SyncLogger SyncLogger) Verbose(message string, sprintf ...interface{}) {
	if Verbose {
		if len(sprintf) > 0 {
			message = fmt.Sprintf(message, sprintf...)
		}

		SyncLogger.Println(message)
	}
}

func (SyncLogger SyncLogger) Main(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	SyncLogger.Println(":: " + message)
}

func (SyncLogger SyncLogger) Step(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	SyncLogger.Println("   -> " + message)
}

func (SyncLogger SyncLogger) FatalExit(exitCode int, message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	SyncLogger.Fatal(message)
	os.Exit(exitCode)
}


// Log error object as message
func (SyncLogger SyncLogger) FatalErrorExit(exitCode int, err error) {

	if CommandName != "" {
		cmdline := fmt.Sprintf("%s %s", CommandName, strings.Join(os.Args[1:], " "))
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Command: %s", cmdline))
	}

	fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: %s", err))

	os.Exit(exitCode)
}
