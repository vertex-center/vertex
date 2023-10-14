package types

type Update struct {
	Baseline Baseline `json:"baseline"`
}

type Updater interface {
	CurrentVersion() (string, error)
	Install(version string) error
	IsInstalled() bool
	ID() string
}
