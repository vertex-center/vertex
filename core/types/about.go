package types

type About struct {
	Version string `json:"version" fake:"{appversion}"`
	Commit  string `json:"commit"  fake:"{commit}"`
	Date    string `json:"date"    fake:"{date}"`
	OS      string `json:"os"      fake:"linux"`
	Arch    string `json:"arch"    fake:"amd64"`
}
