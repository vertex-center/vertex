package service

import (
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type envService struct {
	adapter port.EnvAdapter
}

func NewEnvService(adapter port.EnvAdapter) port.EnvService {
	return &envService{
		adapter: adapter,
	}
}

func (s *envService) Save(inst *types.Container, env types.ContainerEnvVariables) error {
	inst.Env = env
	return s.adapter.Save(inst.UUID, env)
}

func (s *envService) Load(inst *types.Container) error {
	env, err := s.adapter.Load(inst.UUID)
	if err != nil {
		return err
	}
	inst.Env = env
	return nil
}
