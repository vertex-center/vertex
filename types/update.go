package types

import "errors"

var (
	ErrAlreadyUpdating = errors.New("an update is already in progress, cannot start another")
)

type Update struct {
	Baseline Baseline `json:"baseline"`
}

type Updater interface {
	CurrentVersion() (string, error)
	Install(version string) error
	IsInstalled() bool
	ID() string
}
