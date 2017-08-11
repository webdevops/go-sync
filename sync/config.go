package sync

import (
	"regexp"
	"sync"
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
}

type Connection struct {
	Type string
	Hostname string
	User string
	Password string
	Docker string
	WorkDir string
}

type Database struct {
	Type string
	Schema string
	Hostname string
	Port string
	User string
	Password string

	Filter Filter
	Local struct {
		Type string
		Schema string
		Hostname string
		Port string
		User string
		Password string

		Connection Connection
	}
	Options struct {
		ClearDatabase bool `yaml:"clear-database"`
	}

	// local cache
	cacheRemoteTableList []string
	cacheLocalTableList []string

	remoteConnection Connection
}

type Execution struct {
	Type string
	Command YamlStringArray
	Workdir string
	When string
}

type Server struct {
	Path string
	Connection Connection
	Filesystem []Filesystem
	Database []Database
	Exec []Execution
}

type SyncConfig struct {
	Sync map[string]Server
	Deploy map[string]Server
}
