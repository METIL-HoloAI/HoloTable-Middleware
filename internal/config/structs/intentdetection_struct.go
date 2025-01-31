package structs

type IntentDetection struct {
	Endpoint string            `yaml:"endpoint"`
	Method   string            `yaml:"method"`
	Headers  map[string]string `yaml:"headers"`
	// Body               map[string]string           `yaml:"body"` // old one
	Payload       PayloadConfig `yaml:"payload"`
	InitialPrompt string        `yaml:"initialPrompt"`
}

type PayloadConfig struct {
	Model    string            `yaml:"model"`
	Messages []MessageTemplate `yaml:"messages"`
}

type MessageTemplate struct {
	Role    string `yaml:"role"`
	Content string `yaml:"content"`
}
