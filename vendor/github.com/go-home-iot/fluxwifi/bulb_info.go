package fluxwifi

// BulbInfo contains the information returned from scanning the local
// network for Flux WIFI bulbs
type BulbInfo struct {
	IP    string
	ID    string
	Model string
}
