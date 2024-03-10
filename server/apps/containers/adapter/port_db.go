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

type portDBAdapter struct {
	db storage.DB
}

func NewPortDBAdapter(db storage.DB) port.PortAdapter {
	return &portDBAdapter{db}
}

func (a *portDBAdapter) GetPorts(ctx context.Context, filters types.PortFilters) (types.Ports, error) {
	var ports types.Ports
	query := `SELECT * FROM ports`
	var args []interface{}
	if filters.ContainerID != nil {
		query += ` WHERE container_id = $1`
		args = append(args, *filters.ContainerID)
	}
	err := a.db.Select(&ports, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return ports, nil
	}
	return ports, err
}

func (a *portDBAdapter) CreatePort(ctx context.Context, port types.Port) error {
	_, err := a.db.NamedExec(`
		INSERT INTO ports (id, container_id, internal_port, external_port)
		VALUES (:id, :container_id, :internal_port, :external_port)
	`, port)
	return err
}

func (a *portDBAdapter) DeletePort(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(`
        DELETE FROM ports
        WHERE id = $1
    `, id)
	return err
}

func (a *portDBAdapter) DeletePorts(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(`
		DELETE FROM ports
		WHERE container_id = $1
	`, id)
	return err
}

func (a *portDBAdapter) UpdatePortByID(ctx context.Context, port types.Port) error {
	_, err := a.db.NamedExec(`
        UPDATE ports
        SET internal_port = :internal_port, external_port = :external_port
        WHERE id = :id
    `, port)
	return err
}
