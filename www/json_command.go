package www

type jsonCommand struct {
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}
