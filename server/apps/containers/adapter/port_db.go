package adapter

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/storage"
)

type portDBAdapter struct {
	db storage.DB
}

func NewPortDBAdapter(db storage.DB) port.PortAdapter {
	return &portDBAdapter{db}
}

func (a *portDBAdapter) GetContainerPorts(ctx context.Context, id uuid.UUID) (types.Ports, error) {
	var ports types.Ports
	err := a.db.Select(&ports, `
		SELECT * FROM ports
		WHERE container_id = $1
	`, id)
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

func (a *portDBAdapter) DeleteContainerPorts(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(`
		DELETE FROM ports
		WHERE container_id = $1
	`, id)
	return err
}

func (a *portDBAdapter) UpdateContainerPortByID(ctx context.Context, port types.Port) error {
	_, err := a.db.NamedExec(`
        UPDATE ports
        SET internal_port = :internal_port, external_port = :external_port
        WHERE id = :id
    `, port)
	return err
}
