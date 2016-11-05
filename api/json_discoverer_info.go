package api

type jsonDiscovererInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	PreScanInfo string `json:"preScanInfo"`
}
