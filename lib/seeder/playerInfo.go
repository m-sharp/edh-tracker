package seeder

// PlayerInfo holds per-player role configuration from playerInfos.json.
type PlayerInfo struct {
	Name     string `json:"name"`
	UserRole string `json:"userRole"`
	PodRole  string `json:"podRole"`
}
