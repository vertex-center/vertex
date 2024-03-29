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

type capDBAdapter struct {
	db storage.DB
}

func NewCapDBAdapter(db storage.DB) port.CapAdapter {
	return &capDBAdapter{db}
}

func (a *capDBAdapter) GetContainerCaps(ctx context.Context, id uuid.UUID) (types.Capabilities, error) {
	var caps types.Capabilities
	err := a.db.Select(&caps, `
		SELECT * FROM capabilities
		WHERE container_id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return caps, nil
	}
	return caps, err
}

func (a *capDBAdapter) CreateCap(ctx context.Context, c types.Capability) error {
	_, err := a.db.NamedExec(`
		INSERT INTO capabilities (id, container_id, name)
		VALUES (:id, :container_id, :name)
	`, c)
	return err
}

func (a *capDBAdapter) DeleteContainerCaps(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(`
		DELETE FROM capabilities
		WHERE container_id = $1
	`, id)
	return err
}
