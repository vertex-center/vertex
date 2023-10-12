package types

import (
	"io"

	"github.com/vertex-center/vertex/types"
)

type InstanceRunnerAdapterPort interface {
	Delete(inst *Instance) error
	Start(inst *Instance, setStatus func(status string)) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
	Stop(inst *Instance) error
	Info(inst Instance) (map[string]any, error)
	WaitCondition(inst *Instance, cond types.WaitContainerCondition) error

	CheckForUpdates(inst *Instance) error
	HasUpdateAvailable(inst Instance) (bool, error)
	GetAllVersions(inst Instance) ([]string, error)
}
