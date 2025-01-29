package structs

type IntentDetection struct {
	Endpoint      string            `yaml:"endpoint"`
	Method        string            `yaml:"method"`
	InitialPrompt string            `yaml:"initialPrompt"`
	Headers       map[string]string `yaml:"headers"`
	Body          map[string]string `yaml:"body"`
}
