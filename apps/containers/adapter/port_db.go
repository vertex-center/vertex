package adapter

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
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

func (a *portDBAdapter) GetPorts(ctx context.Context, id types.ContainerID) (types.Ports, error) {
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

func (a *portDBAdapter) CreatePorts(ctx context.Context, ports types.Ports) error {
	for _, p := range ports {
		_, err := a.db.NamedExec(`
			INSERT INTO ports (container_id, internal_port, external_port)
			VALUES (:container_id, :internal_port, :external_port)
		`, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *portDBAdapter) DeletePorts(ctx context.Context, id types.ContainerID) error {
	_, err := a.db.Exec(`
		DELETE FROM ports
		WHERE container_id = $1
	`, id)
	return err
}
