package types

type About struct {
	Version string `json:"version" example:"1.2.3"`
	Commit  string `json:"commit"  example:"cd4ba2876f45775287f426c13adb1868f7c96222"`
	Date    string `json:"date"    example:"2006-01-02T15:04:05Z07:00"`
	OS      string `json:"os"      example:"linux"`
	Arch    string `json:"arch"    example:"amd64"`
}
