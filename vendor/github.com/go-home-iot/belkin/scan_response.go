package belkin

// ScanResponse contains information from a device that responded to a scan response
type ScanResponse struct {
	MaxAge     int
	SearchType string
	DeviceID   string
	USN        string
	Location   string
	Server     string
	URN        string
}
