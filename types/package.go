package types

type Package struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Homepage       string            `json:"homepage"`
	License        string            `json:"license"`
	Check          string            `json:"check"`
	InstallPackage map[string]string `json:"install"`
}
