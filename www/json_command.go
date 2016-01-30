package www

type jsonCommand struct {
	Type string `json:"type"`

	// ClientID is an ID assigned on the client, if this is a new command that
	// does not exist on the server.
	ClientID string `json:"clientId"`

	Attributes map[string]interface{} `json:"attributes"`
}
