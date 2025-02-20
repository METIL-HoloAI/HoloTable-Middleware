package structs

import _ "gopkg.in/yaml.v3"

type SpeechToTextSettings struct {
	WebsocketURL string `yaml:"WebsocketURL"`
	Keyword      string `yaml:"Keyword"`
}
