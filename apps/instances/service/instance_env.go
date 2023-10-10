package service

import (
	"github.com/vertex-center/vertex/apps/instances/types"
)

type InstanceEnvService struct {
	adapter types.InstanceEnvAdapterPort
}

func NewInstanceEnvService(adapter types.InstanceEnvAdapterPort) *InstanceEnvService {
	return &InstanceEnvService{
		adapter: adapter,
	}
}

func (s *InstanceEnvService) Save(inst *types.Instance, env types.InstanceEnvVariables) error {
	inst.Env = env
	return s.adapter.Save(inst.UUID, env)
}

func (s *InstanceEnvService) Load(inst *types.Instance) error {
	env, err := s.adapter.Load(inst.UUID)
	if err != nil {
		return err
	}
	inst.Env = env
	return nil
}

func (s *InstanceEnvService) OnEvent(e interface{}) {
	// TODO: Useless
}
