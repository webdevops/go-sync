package sync

import (
	"regexp"
	"sync"
)

var waitGroup sync.WaitGroup

type Filter struct {
	// Exclude as strings (regexp)
	Exclude []string `yaml:"exclude"`
	// compiled regexp excludes
	excludeRegexp []*regexp.Regexp

	// Includes as strings (regexp)
	Include []string `yaml:"include"`
	// compiled regexp includes
	includeRegexp []*regexp.Regexp
}

type Filesystem struct {
	// Remove path
	Path string `yaml:"path"`
	// Local path (optional)
	Local string `yaml:"local"`
	// Filter
	Filter Filter `yaml:"filter"`
	// Connection for filesystem sync (optional, default is Server connection)
	Connection *YamlCommandBuilderConnection `yaml:"connection"`
	Options struct {
		// Generate stubs (small example files) instead of fetching files from remote
		GenerateStubs bool `yaml:"generate-stubs"`
	} `yaml:"options"`
}

type DatabaseOptions struct {
	// Clear database with DROP/CREATE before sync
	ClearDatabase bool `yaml:"clear-database"`
	// Arguments for mysqldump command
	Mysqldump *YamlStringArray `yaml:"mysqldump"`
	// Arguments for mysql command
	Mysql *YamlStringArray `yaml:"mysql"`
	// Arguments for pgdump command
	Pgdump *YamlStringArray `yaml:"pgdump"`
	// Arguments for psql command
	Psql *YamlStringArray `yaml:"psql"`
}

type EnvironmentVar struct {
	// Name of variable
	Name string `yaml:"name"`
	// Value of variable
	Value string `yaml:"value"`
}

type Database struct {
	// Type of database (either mysql or postgres)
	Type string `yaml:"type"`
	// Database name on remote database server
	Db string `yaml:"database"`
	// Hostname of remote database server
	Hostname string `yaml:"hostname"`
	// Port of remote database server
	Port string `yaml:"port"`
	// Username of remote database server
	User string `yaml:"user"`
	// Password of remote database server
	Password string `yaml:"password"`

	// Table filter
	Filter Filter `yaml:"filter"`
	// Connection for database sync (optional, default is Server connection)
	Connection *YamlCommandBuilderConnection `yaml:"connection"`
	// Database options
	Options DatabaseOptions `yaml:"options"`

	Local struct {
		// Database name on local database server
		Db string `yaml:"database"`
		// Hostname of local database server
		Hostname string `yaml:"hostname"`
		// Port of local database server
		Port string `yaml:"port"`
		// Username of local database server
		User string `yaml:"user"`
		// Password of local database server
		Password string `yaml:"password"`

		// Connection for database sync (optional, default is empty)
		Connection *YamlCommandBuilderConnection `yaml:"connection"`

		// Database options
		Options DatabaseOptions `yaml:"options"`
	} `yaml:"local"`

	// local cache for remote table list
	cacheRemoteTableList []string
	// local cache for local table list
	cacheLocalTableList []string
}

type Execution struct {
	// Type of execution (remote or local)
	Type string `yaml:"type"`
	// Command as string or as elements
	Command YamlStringArray `yaml:"command"`
	// Workdir for execution
	Workdir string `yaml:"workdir"`

	// Environment variables
	Environment []EnvironmentVar `yaml:"environment"`

	// Execution options
	Options struct {
	} `yaml:"options"`
}

type Server struct {
	// General working path (for filesystem syncs)
	Path string `yaml:"path"`
	// General connection (default for all remote connections)
	Connection *YamlCommandBuilderConnection `yaml:"connection"`
	// Filesystem sync list
	Filesystem []Filesystem `yaml:"filesystem"`
	// Database sync list
	Database []Database `yaml:"database"`
	// Startup execution list (executed before sync)
	ExecStartup []Execution `yaml:"exec-startup"`
	// Finish execution list (executed after sync)
	ExecFinish []Execution `yaml:"exec-finish"`

	runConfiguration *RunConfiguration
}

type SyncConfig struct {
	// Sync (remote -> local) configurations
	Sync map[string]Server `yaml:"sync"`
	// Deploy (local -> remote) configurations
	Deploy map[string]Server `yaml:"deploy"`
}

type RunConfiguration struct {
	// Enable database sync
	Database   bool
	// Enable filesystem sync
	Filesystem bool
	// Enable exec runner
	Exec       bool
}
