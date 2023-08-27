package types

type About struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`

	OS   string `json:"os"`
	Arch string `json:"arch"`
}
