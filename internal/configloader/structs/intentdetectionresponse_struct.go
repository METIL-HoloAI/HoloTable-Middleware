package structs

type IntentDetectionResponse struct {
	ContentType        string            `json:"contentType"`
	Endpoint           string            `json:"endpoint"`
	Method             string            `json:"method"`
	Headers            map[string]string `json:"headers"`
	RequiredParameters map[string]struct {
		Description string   `json:"description"`
		Options     []string `json:"options"`
	} `json:"requiredParameters"`
	OptionalParameters map[string]struct {
		Default     interface{}   `json:"default"` // Use interface{} for mixed types like string/int
		Description string        `json:"description"`
		Options     []interface{} `json:"options"` // Use interface{} for mixed types like int/string
	} `json:"optionalParameters"`
}
