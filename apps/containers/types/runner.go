package types

import (
	"io"

	"github.com/vertex-center/vertex/types"
)

type ContainerRunnerAdapterPort interface {
	Delete(inst *Container) error
	Start(inst *Container, setStatus func(status string)) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
	Stop(inst *Container) error
	Info(inst Container) (map[string]any, error)
	WaitCondition(inst *Container, cond types.WaitContainerCondition) error

	CheckForUpdates(inst *Container) error
	HasUpdateAvailable(inst Container) (bool, error)
	GetAllVersions(inst Container) ([]string, error)
}
