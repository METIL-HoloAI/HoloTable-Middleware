package structs

// Represents a single API call step
type Step struct {
	Name                 string                 `yaml:"name"`                            // Step identifier
	Method               string                 `yaml:"method"`                          // HTTP method (GET, POST, etc.)
	URL                  string                 `yaml:"url"`                             // API endpoint (with placeholders)
	Headers              map[string]string      `yaml:"headers,omitempty"`               // Optional headers for the request
	Body                 map[string]interface{} `yaml:"body,omitempty"`                  // Request payload (now fully dynamic)
	ResponsePlaceholders map[string]interface{} `yaml:"response_placeholders,omitempty"` // Mapping placeholders -> response field names
	Poll                 map[string]interface{} `yaml:"poll,omitempty"`                  // Fully dynamic polling config
	ContentExtraction    map[string]interface{} `yaml:"content_extraction,omitempty"`
}

// Represents a full API workflow (multiple steps)
type Workflow struct {
	Steps []Step `yaml:"steps"`
}

// Maps content type (e.g., "imagegen", "3dgen") to its workflow
type WorkflowCollection map[string]Workflow
