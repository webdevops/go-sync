package sync

import (
	"regexp"
	"sync"
	"github.com/webdevops/go-shell/commandbuilder"
)

var waitGroup sync.WaitGroup

type Filter struct {
	Exclude []string
	excludeRegexp []*regexp.Regexp
	Include []string
	includeRegexp []*regexp.Regexp
}

type Filesystem struct {
	Path string
	Local string
	Filter Filter
	Connection commandbuilder.Connection
	Options struct {
		GenerateStubs bool `yaml:"generate-stubs"`
	}
}

type DatabaseOptions struct {
	ClearDatabase bool `yaml:"clear-database"`
	Mysqldump string
	Mysql string
}

type Database struct {
	Type string
	Schema string
	Hostname string
	Port string
	User string
	Password string

	Filter Filter
	Connection commandbuilder.Connection

	Local struct {
		Type string
		Schema string
		Hostname string
		Port string
		User string
		Password string

		Connection commandbuilder.Connection
		Options DatabaseOptions
	}
	Options DatabaseOptions

	// local cache
	cacheRemoteTableList []string
	cacheLocalTableList []string
}

type Execution struct {
	Type string
	Command YamlStringArray
	Workdir string
}

type Server struct {
	Path string
	Connection commandbuilder.Connection
	Filesystem []Filesystem
	Database []Database
	ExecStartup []Execution `yaml:"exec-startup"`
	ExecFinish []Execution `yaml:"exec-finish"`
}

type SyncConfig struct {
	Sync map[string]Server
	Deploy map[string]Server
}
