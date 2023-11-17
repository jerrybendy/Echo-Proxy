package service

type Setting struct {
	HttpPort  int16 `json:"httpPort"`
	HttpsPort int64 `json:"httpsPort"`
}

func (s *Setting) OpenConfigFolder() {
	configDir := getConfigPath()
	err := openFolder(configDir)
	if err != nil {
		emitErrorToFrontend(err.Error())
	}
}

func (s *Setting) GetSettings() *Setting {
	if s.HttpPort == 0 {
		s.HttpPort = 80
	}
	if s.HttpsPort == 0 {
		s.HttpsPort = 443
	}
	return s
}

func (s *Setting) SaveSettings(set Setting) bool {
	config.Setting = set
	saveConfig()
	return true
}
