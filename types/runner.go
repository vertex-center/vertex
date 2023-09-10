package types

type RunnerAdapterPort interface {
	Delete(instance *Instance) error
	Start(instance *Instance, onLog func(msg string), onErr func(msg string), setStatus func(status string)) error
	Stop(instance *Instance) error
	Info(instance Instance) (map[string]any, error)

	CheckForUpdates(instance *Instance) error
	HasUpdateAvailable(instance Instance) (bool, error)
}
