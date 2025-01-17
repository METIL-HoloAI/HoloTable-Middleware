package structs

type APIConfig struct {
	Endpoint           string
	Method             string
	Headers            map[string]string
	RequiredParameters map[string]string
	OptionalParameters map[string]string
}
