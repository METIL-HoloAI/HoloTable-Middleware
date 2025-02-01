package structs

type IntentDetection struct {
	Endpoint      string                 `yaml:"endpoint"`
	Method        string                 `yaml:"method"`
	Headers       map[string]string      `yaml:"headers"`
	Payload       map[string]interface{} `yaml:"payload"`
	InitialPrompt string                 `yaml:"initialPrompt"`
}
