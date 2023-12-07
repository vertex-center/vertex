package service

import (
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type containerEnvService struct {
	adapter port.ContainerEnvAdapter
}

func NewContainerEnvService(adapter port.ContainerEnvAdapter) port.ContainerEnvService {
	return &containerEnvService{
		adapter: adapter,
	}
}

func (s *containerEnvService) Save(inst *types.Container, env types.ContainerEnvVariables) error {
	inst.Env = env
	return s.adapter.Save(inst.UUID, env)
}

func (s *containerEnvService) Load(inst *types.Container) error {
	env, err := s.adapter.Load(inst.UUID)
	if err != nil {
		return err
	}
	inst.Env = env
	return nil
}
