package logger

import (
	"log"
	"os"
	"fmt"
	"strings"
)

const (
	LogPrefix = ""
	prefixMain = ":: "
	prefixSub  = "   -> "
	prefixCmd  = "      $ "
	prefixErr  = " [ERROR] "
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

	SyncLogger.Println(prefixMain + message)
}

func (SyncLogger SyncLogger) Step(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	SyncLogger.Println(prefixSub + message)
}


func (SyncLogger SyncLogger) Command(message string) {
	SyncLogger.Println(prefixCmd + message)
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

	fmt.Fprintln(os.Stderr, fmt.Sprintf("%s %s", prefixErr, err))

	os.Exit(exitCode)
}
