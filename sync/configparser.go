package sync

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	"errors"
	"strings"
)

const globalYamlHeader = `
---
databaseExcludeTYPO3:
  - ^cachingframework_.*
  - ^cf_.*"
  - ^cache_.*
  - ^index_.*
  - ^sys_log$
  - ^sys_history$
  - ^sys_registry$
  - ^tx_extbase_cache.*
  - ^tx_extensionmanager_domain_model_extension.*
  - ^zzz_deleted_.*
---
`

func NewConfigParser(file string) (config *SyncConfig) {
	ymlData, err := ioutil.ReadFile(file)
	if err != nil {
		Logger.FatalErrorExit(1, err)
	}

	// ymlDataCombined := globalYamlHeader + fmt.Sprintf("%s", ymlData)
	ymlDataCombined := fmt.Sprintf("%s", ymlData)

	if err := yaml.Unmarshal([]byte(ymlDataCombined), &config); err != nil {
		Logger.FatalErrorExit(1, err)
	}

	return
}

func (config *SyncConfig) GetSyncServer(serverName string) (Server, error) {
	if val, ok := config.Sync[serverName]; ok {
		return val, nil
	} else {
		return Server{}, errors.New(fmt.Sprintf("Server name %s doesn't exists", serverName))
	}
}

func (config *SyncConfig) GetDeployServer(serverName string) (Server, error) {
	if val, ok := config.Deploy[serverName]; ok {
		return val, nil
	} else {
		return Server{}, errors.New(fmt.Sprintf("Server name %s doesn't exists", serverName))
	}
}

func (config *SyncConfig) GetServerList(confType string) (list []string) {
	switch confType {
	case "sync":
		for key := range config.Sync {
			list = append(list, key)
		}
	case "deploy":
		for key := range config.Deploy {
			list = append(list, key)
		}
	}

	return
}

// List all possible server configurations
func (config *SyncConfig) ListServer() (list map[string][]string) {
	if len(config.Sync) > 0 {
		list["Sync"] = make([]string, len(config.Sync)-1)
		for key := range config.Sync {
			list["Sync"] = append(list["Sync"], key)
		}
	}

	if len(config.Deploy) > 0 {
		list["Deploy"] = make([]string, len(config.Deploy)-1)
		for key := range config.Deploy {
			list["Deploy"] = append(list["Deploy"], key)
		}
	}

	return
}

// Show all possible server configurations
// in an human readable style
func (config *SyncConfig) ShowConfiguration() {
	serverList := config.ListServer()

	for area, keyList := range serverList {
		fmt.Println()
		fmt.Println(area)
		fmt.Println(strings.Repeat("=", len(area)))

		for _, serverKey := range keyList {
			fmt.Println(fmt.Sprintf(" -> %s ", serverKey))
		}
	}


}
