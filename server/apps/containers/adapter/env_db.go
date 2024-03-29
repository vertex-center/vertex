package adapter

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/containers/core/port"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
	"github.com/vertex-center/vertex/server/common/storage"
)

type envDBAdapter struct {
	db storage.DB
}

func NewEnvDBAdapter(db storage.DB) port.EnvAdapter {
	return &envDBAdapter{db}
}

func (a *envDBAdapter) GetEnvs(ctx context.Context, filters types.EnvVariableFilters) ([]types.EnvVariable, error) {
	var env []types.EnvVariable
	query := `SELECT * FROM env_variables`
	var args []interface{}
	if filters.ContainerID != nil {
		query += ` WHERE container_id = $1`
		args = append(args, *filters.ContainerID)
	}
	query += ` ORDER BY name`
	err := a.db.Select(&env, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return env, nil
	}
	return env, err
}

func (a *envDBAdapter) CreateEnv(ctx context.Context, v types.EnvVariable) error {
	_, err := a.db.NamedExec(`
		INSERT INTO env_variables (id, container_id, type, name, display_name, value, default_value, description, secret)
		VALUES (:id, :container_id, :type, :name, :display_name, :value, :default_value, :description, :secret)
	`, v)
	return err
}

func (a *envDBAdapter) DeleteEnv(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(`
        DELETE FROM env_variables
        WHERE id = $1
    `, id)
	return err
}

func (a *envDBAdapter) DeleteEnvs(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(`
		DELETE FROM env_variables
		WHERE container_id = $1
	`, id)
	return err
}

func (a *envDBAdapter) UpdateEnvByID(ctx context.Context, v types.EnvVariable) error {
	_, err := a.db.Exec(`
		UPDATE env_variables
		SET name = $1, value = $2
		WHERE id = $3
	`, v.Name, v.Value, v.ID)
	return err
}

func (a *envDBAdapter) UpdateEnvByName(ctx context.Context, v types.EnvVariable) error {
	_, err := a.db.Exec(`
		UPDATE env_variables
		SET value = $1
		WHERE container_id = $2 AND name = $3
	`, v.Value, v.ContainerID, v.Name)
	return err
}
