package api

import "strings"

type jsonZone struct {
	Address     string `json:"address"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DeviceID    string `json:"deviceId"`
	Type        string `json:"type"`
	Output      string `json:"output"`
}
type zones []jsonZone

func (slice zones) Len() int {
	return len(slice)
}
func (slice zones) Less(i, j int) bool {
	return strings.ToLower(slice[i].Name) < strings.ToLower(slice[j].Name)
}
func (slice zones) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
