package service

import (
	"context"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type envService struct {
	env port.EnvAdapter
}

func NewEnvService(env port.EnvAdapter) port.EnvService {
	return &envService{env}
}

func (s *envService) GetEnvs(ctx context.Context, filters types.EnvVariableFilters) ([]types.EnvVariable, error) {
	return s.env.GetEnvs(ctx, filters)
}

func (s *envService) PatchEnv(ctx context.Context, env types.EnvVariable) error {
	err := env.Validate()
	if err != nil {
		return err
	}
	return s.env.UpdateEnvByID(ctx, env)
}

func (s *envService) DeleteEnv(ctx context.Context, id uuid.UUID) error {
	return s.env.DeleteEnv(ctx, id)
}

func (s *envService) CreateEnv(ctx context.Context, env types.EnvVariable) error {
	env.ID = uuid.New()
	err := env.Validate()
	if err != nil {
		return err
	}
	return s.env.CreateEnv(ctx, env)
}
