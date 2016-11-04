package api

type jsonDiscovererInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Type can either be ScanDevices|FromString
	Type string `json:"type"`

	//TODO: Expose image asset
}
