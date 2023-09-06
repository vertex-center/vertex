package types

// Dependency describes main vertex programs that Vertex needs to run, like Vertex,
// Vertex Services, Vertex Web UI...
type Dependency interface {
	// CheckForUpdate will check if the dependency has an update available. If
	// true, it returns the Update, or nil otherwise.
	CheckForUpdate() (*Update, error)

	// InstallUpdate will install the previously fetched update.
	InstallUpdate() error

	// GetID returns the ID of the dependency that can be used to
	// identify it from the client.
	GetID() string

	// GetPath returns the path of the dependency.
	GetPath() string
}
