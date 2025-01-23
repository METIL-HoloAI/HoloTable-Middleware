package structs

type ChatConfig struct {
    Endpoint           string                      `yaml:"endpoint"`
    Method             string                      `yaml:"method"`
    Headers            map[string]string           `yaml:"headers"`
    RequiredParameters map[string]ParameterDetails `yaml:"requiredParameters"`
    OptionalParameters map[string]ParameterDetails `yaml:"optionalParameters"`
}

type ParameterDetails struct {
    Description string        `yaml:"description"`
    Default     interface{}   `yaml:"default,omitempty"`
    Options     []interface{} `yaml:"options,omitempty"`
}