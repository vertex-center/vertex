package types

const (
	PmNone   = "sources"
	PmAptGet = "apt-get"
	PmBrew   = "brew"
	PmNpm    = "npm"
	PmPacman = "pacman"
	PmSnap   = "snap"
)

type Package struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Homepage       string            `json:"homepage"`
	License        string            `json:"license"`
	Check          string            `json:"check"`
	InstallPackage map[string]string `json:"install"`
	Installed      *bool             `json:"installed,omitempty"`
}

type PackageAdapterPort interface {
	GetByID(id string) (Package, error)
	GetPath(id string) string

	// Reload the adapter
	Reload() error
}
