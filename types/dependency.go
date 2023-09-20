package types

import "time"

// DependencyUpdater describes main vertex programs that Vertex needs to run, like Vertex,
// Vertex Services, Vertex Web UI...
type DependencyUpdater interface {
	// CheckForUpdate will check if the dependency has an update available. If
	// true, it returns the Dependency, or nil otherwise.
	CheckForUpdate() (*DependencyUpdate, error)

	// InstallUpdate will install the previously fetched update.
	InstallUpdate() error

	// GetID returns the ID of the dependency that can be used to
	// identify it from the client.
	GetID() string

	// GetPath returns the path of the dependency.
	GetPath() string
}

type Dependency struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`

	Update  *DependencyUpdate `json:"update,omitempty"`
	Updater DependencyUpdater `json:"-"`
}

type DependencyUpdate struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	NeedsRestart   bool   `json:"needs_restart,omitempty"`
}

type Dependencies struct {
	LastUpdatesCheck *time.Time    `json:"last_updates_check"`
	Items            []*Dependency `json:"items"`
}
