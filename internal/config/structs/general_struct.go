package structs

import _ "gopkg.in/yaml.v3"

type GeneralSettings struct {
	Listener      string `yaml:"listener"`
	DataDir       string `yaml:"dataDir"`
	OpenWebsocket bool   `yaml:"openWebsocket"`
}
