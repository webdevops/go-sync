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



func NewConfigParser(file string) *SyncConfig {
	config := SyncConfig{}

	ymlData, err := ioutil.ReadFile(file)
	if err != nil {
		Logger.FatalErrorExit(1, err)
	}

	// ymlDataCombined := globalYamlHeader + fmt.Sprintf("%s", ymlData)
	ymlDataCombined := fmt.Sprintf("%s", ymlData)

	if err := yaml.Unmarshal([]byte(ymlDataCombined), &config); err != nil {
		Logger.FatalErrorExit(1, err)
	}

	return &config
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

func (config *SyncConfig) GetServerList(confType string) []string {
	ret := []string{}

	switch confType {
	case "sync":
		for key := range config.Sync {
			ret = append(ret, key)
		}
	case "deploy":
		for key := range config.Deploy {
			ret = append(ret, key)
		}
	}

	return ret
}

func (config *SyncConfig) ListServer() map[string][]string {
	ret := map[string][]string{}

	ret["Sync"] = make([]string, len(config.Sync)-1)
	for key := range config.Sync {
		ret["Sync"] = append(ret["Sync"], key)
	}

	ret["Deploy"] = make([]string, len(config.Deploy)-1)
	for key := range config.Deploy {
		ret["Deploy"] = append(ret["Deploy"], key)
	}

	return ret
}

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
