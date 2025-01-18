package structs

type IntentDetectionResponse struct {
	ContentType        string
	Prompt             string
	RequiredParamaters map[string]string
	OptionalParameters map[string]string
}
