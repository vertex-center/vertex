package adapter

import (
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/storage"
)

type envDBAdapter struct {
	db storage.DB
}

func NewEnvDBAdapter(db storage.DB) port.EnvAdapter {
	return &envDBAdapter{db}
}

func (a *envDBAdapter) GetEnv(uuid types.ContainerID) (types.EnvVariables, error) {
	var env types.EnvVariables
	err := a.db.Select(&env, `
		SELECT * FROM env_variables
		WHERE container_id = $1
	`, uuid)
	return env, err
}

func (a *envDBAdapter) SaveEnv(uuid types.ContainerID, variables types.EnvVariables) error {
	// TODO: Implement
	return nil
}
