package service

import (
	"context"
	"encoding/json"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"localProxy/utils"
	"os"
)

const AppName = "Echo Proxy"

var configFilePath string

var config *configFile

type configFile struct {
	Version string        `json:"version"`
	Hosts   []*HostConfig `json:"hosts"`
	Setting Setting       `json:"setting"`
}

var runtimeCtx context.Context

func Init(ctx context.Context) {
	runtimeCtx = ctx

	configDir := getConfigPath()
	if !utils.FileExists(configDir) {
		err := os.Mkdir(configDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	if configFilePath == "" {
		configFilePath = configDir + string(os.PathSeparator) + "config.json"
	}

	parseConfig()
}

func getConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	configDir = configDir + string(os.PathSeparator) + AppName
	return configDir
}

func parseConfig() {
	if !utils.FileExists(configFilePath) {
		config = &configFile{}
		config.Hosts = make([]*HostConfig, 0)
	} else {
		fileContent, err := os.ReadFile(configFilePath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(fileContent, &config)
		if err != nil {
			panic(err)
		}
	}
}

func saveConfig() {
	go func() {
		content, err := json.Marshal(config)
		if err != nil {
			emitErrorToFrontend("Encode config failed, " + err.Error())
			return
		}
		err = os.WriteFile(configFilePath, content, 0755)
		if err != nil {
			emitErrorToFrontend("Save config failed, " + err.Error())
			return
		}
		runtime.EventsEmit(runtimeCtx, "hostsChange")
	}()
}

func emitErrorToFrontend(errorMessage string) {
	runtime.EventsEmit(runtimeCtx, "error", errorMessage)
}

func emitWarningToFrontend(errorMessage string) {
	runtime.EventsEmit(runtimeCtx, "warning", errorMessage)
}
