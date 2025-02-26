package structs

import _ "gopkg.in/yaml.v3"

type SpeechToTextSettings struct {
	LiveTranscription struct {
		WebsocketURL string `yaml:"websocketURL"`
		Keyword      string `yaml:"keyword"`
	} `yaml:"liveTranscription"`

	SpeechToTextAPI struct {
		Endpoint     string                 `yaml:"endpoint"`
		Method       string                 `yaml:"method"`
		Headers      map[string]string      `yaml:"headers"`
		Payload      map[string]interface{} `yaml:"payload"`
		ResponsePath string                 `yaml:"responsePath"`
	} `yaml:"speechToTextAPI"`
}
