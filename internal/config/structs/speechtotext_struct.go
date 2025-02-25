package structs

import _ "gopkg.in/yaml.v3"

type SpeechToTextSettings struct {
	WebsocketURL string                 `yaml:"WebsocketURL"`
	Keyword      string                 `yaml:"Keyword"`
	Endpoint     string                 `yaml:"endpoint"`
	Method       string                 `yaml:"method"`
	Headers      map[string]string      `yaml:"headers"`
	Payload      map[string]interface{} `yaml:"payload"`
	ResponsePath string                 `yaml:"responsePath"`
}
