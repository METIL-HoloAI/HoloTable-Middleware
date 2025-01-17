package structs

type APIConfig struct {
	Endpoint           string            `yaml:"endpoint"`
	Method             string            `yaml:"method"`
	Headers            map[string]string `yaml:"headers"`
	RequiredParameters map[string]string `yaml:"requiredParameters"`
	OptionalParameters map[string]string `yaml:"optionalParameters"`
}
