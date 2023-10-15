package service

import (
	"github.com/vertex-center/vertex/apps/containers/core/port"
	types2 "github.com/vertex-center/vertex/apps/containers/core/types"
)

type ContainerEnvService struct {
	adapter port.ContainerEnvAdapter
}

func NewContainerEnvService(adapter port.ContainerEnvAdapter) *ContainerEnvService {
	return &ContainerEnvService{
		adapter: adapter,
	}
}

func (s *ContainerEnvService) Save(inst *types2.Container, env types2.ContainerEnvVariables) error {
	inst.Env = env
	return s.adapter.Save(inst.UUID, env)
}

func (s *ContainerEnvService) Load(inst *types2.Container) error {
	env, err := s.adapter.Load(inst.UUID)
	if err != nil {
		return err
	}
	inst.Env = env
	return nil
}
