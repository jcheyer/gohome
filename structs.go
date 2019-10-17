package gohome

type discoveryInfo struct {
	BaseURL             string `json:"base_url"`
	LocationName        string `json:"location_name"`
	RequiresAPIPassword bool   `json:"requires_api_password"`
	Version             string `json:"version"`
}
