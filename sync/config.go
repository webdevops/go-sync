package sync

import (
	"regexp"
	"sync"
)

var waitGroup sync.WaitGroup

type Filter struct {
	Exclude []string `yaml:"exclude"`
	excludeRegexp []*regexp.Regexp
	Include []string `yaml:"include"`
	includeRegexp []*regexp.Regexp
}

type Filesystem struct {
	Path string `yaml:"path"`
	Local string `yaml:"local"`
	Filter Filter `yaml:"filter"`
	Connection *YamlCommandBuilderConnection `yaml:"connection"`
	Options struct {
		GenerateStubs bool `yaml:"generate-stubs"`
	} `yaml:"options"`
}

type DatabaseOptions struct {
	ClearDatabase bool `yaml:"clear-database"`
	Mysqldump string `yaml:"mysqldump"`
	Mysql string `yaml:"mysql"`
	Pgdump string `yaml:"pgdump"`
	Psql string `yaml:"psql"`
}

type EnvironmentVar struct {
	Name string `yaml:"name"`
	Value string `yaml:"value"`
}

type Database struct {
	Type string `yaml:"type"`
	Db string `yaml:"database"`
	Hostname string `yaml:"hostname"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Password string `yaml:"password"`

	Filter Filter `yaml:"filter"`
	Connection *YamlCommandBuilderConnection `yaml:"connection"`
	Options DatabaseOptions `yaml:"options"`

	Local struct {
		Type string `yaml:"type"`
		Db string `yaml:"database"`
		Hostname string `yaml:"hostname"`
		Port string `yaml:"port"`
		User string `yaml:"user"`
		Password string `yaml:"password"`

		Connection *YamlCommandBuilderConnection `yaml:"connection"`
		Options DatabaseOptions `yaml:"options"`
	} `yaml:"local"`

	// local cache
	cacheRemoteTableList []string
	cacheLocalTableList []string
}

type Execution struct {
	Type string `yaml:"type"`
	Command YamlStringArray `yaml:"command"`
	Workdir string `yaml:"workdir"`
	Environment []EnvironmentVar `yaml:"environment"`
	Options struct {
	} `yaml:"options"`
}

type Server struct {
	Path string `yaml:"path"`
	Connection *YamlCommandBuilderConnection `yaml:"connection"`
	Filesystem []Filesystem `yaml:"filesystem"`
	Database []Database `yaml:"database"`
	ExecStartup []Execution `yaml:"exec-startup"`
	ExecFinish []Execution `yaml:"exec-finish"`

	runConfiguration *RunConfiguration
}

type SyncConfig struct {
	Sync map[string]Server `yaml:"sync"`
	Deploy map[string]Server `yaml:"deploy"`
}

type RunConfiguration struct {
	Database   bool
	Filesystem bool
	Exec       bool
}
