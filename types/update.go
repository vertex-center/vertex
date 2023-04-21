package types

type Update struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	UpToDate       bool   `json:"up_to_date"`
}
