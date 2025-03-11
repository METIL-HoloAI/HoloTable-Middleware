package structs

import _ "gopkg.in/yaml.v3"

type SpeechToTextSettings struct {
	VoskURL string `yaml:"VoskURL"`
	Keyword string `yaml:"keyword"`
}
