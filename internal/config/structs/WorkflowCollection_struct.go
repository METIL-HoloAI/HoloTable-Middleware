package structs

// Represents a single API call step
type Step struct {
	Name        string            `yaml:"name"`
	Method      string            `yaml:"method"`
	URL         string            `yaml:"url"`
	Body        map[string]string `yaml:"body,omitempty"`
	ResponseKey string            `yaml:"response_key,omitempty"`
	Poll        *PollConfig       `yaml:"poll,omitempty"`
}

// Represents polling configuration for async API calls
type PollConfig struct {
	Until    string `yaml:"until"`
	Interval int    `yaml:"interval"`
}

// Represents the overall workflow
type Workflow struct {
	Steps []Step `yaml:"steps"`
}

// Represents API Configuration
type APIConfig struct {
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

// Top-level struct that holds both workflow and API configuration
type APIYaml struct {
	Workflow  Workflow  `yaml:"workflow"`
	APIConfig APIConfig `yaml:",inline"` // Inline for easy access
}

// Maps content type -> APIYaml (Workflow + Config)
type WorkflowCollection map[string]APIYaml
