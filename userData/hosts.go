package userData

type HostConfig struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	ApplyHosts bool   `json:"applyHosts"`
}

type Hosts struct {
}

func (h *Hosts) GetHosts() []*HostConfig {
	return config.Hosts
}

func (h *Hosts) SaveSetting(setting HostConfig) {
	if setting.ID == 0 {
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
