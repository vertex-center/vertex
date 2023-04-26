package types

import "time"

type Update struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	NeedsRestart   bool   `json:"needs_restart"`
}

type Updates struct {
	LastChecked *time.Time `json:"last_checked"`
	Items       []Update   `json:"items"`
}
