package structs

import _ "gopkg.in/yaml.v3"

type GeneralSettings struct {
	Listener   string `yaml:"listener"`
	DataDir    string `yaml:"dataDir"`
	Log_Level  string `yaml:"log_level"`
	Log_Format string `yaml:"log_format"`
}
