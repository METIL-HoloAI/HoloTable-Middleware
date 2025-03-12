package structs

type IntentDetectionResponse struct {
	ContentType        string                 `json:"contentType"`
	RequiredParameters map[string]interface{} `json:"requiredParameters"`
	OptionalParameters map[string]interface{} `json:"optionalParameters"`
}
