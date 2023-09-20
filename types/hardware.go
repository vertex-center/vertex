package types

type Hardware struct {
	OS      string `json:"os"`
	Arch    string `json:"arch"`
	Version string `json:"version"`
	Name    string `json:"name"`
}
