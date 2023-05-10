package types

type RunnerRepository interface {
	Delete(instance *Instance) error
	Start(instance *Instance) error
	Stop(instance *Instance) error
	Info(instance Instance) (map[string]any, error)
}
