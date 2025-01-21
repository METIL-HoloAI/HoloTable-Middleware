package structs

type IntentDetectionResponse struct {
	ContentType        string                 `json:"ContentType"`
	RequiredParameters map[string]string      `json:"requiredParameters"`
	OptionalParameters map[string]interface{} `json:"optionalParameters"`
}
