package types

type Update struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}
