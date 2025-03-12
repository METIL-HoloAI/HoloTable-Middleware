package structs

import _ "gopkg.in/yaml.v3"

type GeneralSettings struct {
	DataDir       string `yaml:"dataDir"`
	OpenWebsocket bool   `yaml:"openWebsocket"`
	Log_Level     string `yaml:"log_level"`
}
