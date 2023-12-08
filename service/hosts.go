package service

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/url"
	"strings"
	"time"
)

type HostConfig struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"` // Domain name
	ApplyHosts  bool         `json:"applyHosts"`
	EnableTLS   bool         `json:"enableTLS"`
	TLSCertFile string       `json:"TLSCertFile"`
	TLSKeyFile  string       `json:"TLSKeyFile"`
	Proxies     []*HostProxy `json:"proxies"`
}

type HostProxy struct {
	ID        int64     `json:"id"`
	MatchType MatchType `json:"matchType"`
	MatchRule string    `json:"matchRule"`
	//MatchParams struct {
	//	CaseInsensitive bool `json:"caseInsensitive,omitempty"`
	//} `json:"matchParams"`
	Target       string `json:"target"`
	ChangeOrigin bool   `json:"changeOrigin"`
	targetUrl    *url.URL
}

type MatchType string

const (
	MatchTypePrefix MatchType = "PREFIX"
	MatchTypeRegexp           = "REGEXP"
	MatchTypeGlob             = "GLOB"
	MatchTypeStatic           = "STATIC"
)

type Hosts struct {
}

func (h *Hosts) GetHosts() []*HostConfig {
	return config.Hosts
}

func (h *Hosts) SaveSetting(setting HostConfig) {
	if setting.ID == 0 {
		setting.ID = time.Now().UnixMilli()
		setting.Name = strings.ToLower(setting.Name)
		config.Hosts = append(config.Hosts, &setting)
	} else {
		for i, val := range config.Hosts {
			if val.ID == setting.ID {
				config.Hosts[i] = &setting
				break
			}
		}
	}
	saveConfig()
}

func (h *Hosts) RemoveHost(id int64) {
	newHosts := make([]*HostConfig, 0, len(config.Hosts))
	for _, item := range config.Hosts {
		if item.ID != id {
			newHosts = append(newHosts, item)
		}
	}
	config.Hosts = newHosts
	saveConfig()
}

func (h *Hosts) OpenFileDialog(title, filterName, filterPattern string) string {
	file, err := runtime.OpenFileDialog(runtimeCtx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{
				DisplayName: filterName,
				Pattern:     filterPattern,
			},
		},
		CanCreateDirectories:       true,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	})
	if err != nil {
		return ""
	}
	return file
}
