package service

import (
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type ContainerEnvService struct {
	adapter port.ContainerEnvAdapter
}

func NewContainerEnvService(adapter port.ContainerEnvAdapter) port.ContainerEnvService {
	return &ContainerEnvService{
		adapter: adapter,
	}
}

func (s *ContainerEnvService) Save(inst *types.Container, env types.ContainerEnvVariables) error {
	inst.Env = env
	return s.adapter.Save(inst.UUID, env)
}

func (s *ContainerEnvService) Load(inst *types.Container) error {
	env, err := s.adapter.Load(inst.UUID)
	if err != nil {
		return err
	}
	inst.Env = env
	return nil
}
