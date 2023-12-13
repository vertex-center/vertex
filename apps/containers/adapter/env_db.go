package adapter

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
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

func (a *envDBAdapter) GetVariable(ctx context.Context, id types.ContainerID) (types.EnvVariables, error) {
	var env types.EnvVariables
	err := a.db.Select(&env, `
		SELECT * FROM env_variables
		WHERE container_id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return env, nil
	}
	return env, err
}

func (a *envDBAdapter) CreateVariables(ctx context.Context, variables types.EnvVariables) error {
	_, err := a.db.NamedExec(`
		INSERT INTO env_variables (container_id, type, name, value, default_value, description)
		VALUES (:container_id, :type, :name, :value, :default_value, :description)
	`, variables)
	return err
}

func (a *envDBAdapter) DeleteVariables(ctx context.Context, id types.ContainerID) error {
	_, err := a.db.Exec(`
		DELETE FROM env_variables
		WHERE container_id = $1
	`, id)
	return err
}