package structs

type APIConfig struct {
	Endpoint           string            `yaml:"endpoint"`
	Method             string            `yaml:"method"`
	Headers            map[string]string `yaml:"headers"`
	RequiredParameters map[string]struct {
		Description string   `yaml:"description"`
		Options     []string `yaml:"options"`
	} `yaml:"requiredParameters"`
	OptionalParameters map[string]struct {
		Default     interface{}   `yaml:"default"` // Mixed types like int/string
		Description string        `yaml:"description"`
		Options     []interface{} `yaml:"options"` // Mixed types like int/string
	} `yaml:"optionalParameters"`
}

/*
type APIConfig struct {
	Endpoint           string            `yaml:"endpoint"`
	Method             string            `yaml:"method"`
	Headers            map[string]string `yaml:"headers"`
	RequiredParameters map[string]string `yaml:"requiredParameters"`
	OptionalParameters map[string]string `yaml:"optionalParameters"`
}
*/
