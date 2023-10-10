package types

import "io"

type InstanceRunnerAdapterPort interface {
	Delete(instance *Instance) error
	Start(instance *Instance, setStatus func(status string)) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
	Stop(instance *Instance) error
	Info(instance Instance) (map[string]any, error)

	CheckForUpdates(instance *Instance) error
	HasUpdateAvailable(instance Instance) (bool, error)
	GetAllVersions(instance Instance) ([]string, error)
}
