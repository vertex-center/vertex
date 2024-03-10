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

type volumeDBAdapter struct {
	db storage.DB
}

func NewVolumeDBAdapter(db storage.DB) port.VolumeAdapter {
	return &volumeDBAdapter{db}
}

func (a *volumeDBAdapter) GetContainerVolumes(ctx context.Context, id uuid.UUID) (types.Volumes, error) {
	var volumes types.Volumes
	err := a.db.Select(&volumes, `
		SELECT * FROM volumes
		WHERE container_id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return volumes, nil
	}
	return volumes, err
}

func (a *volumeDBAdapter) CreateVolume(ctx context.Context, vol types.Volume) error {
	_, err := a.db.NamedExec(`
			INSERT INTO volumes (id, container_id, type, internal_path, external_path)
			VALUES (:id, :container_id, :type, :internal_path, :external_path)
		`, vol)
	return err
}

func (a *volumeDBAdapter) DeleteContainerVolumes(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(`
		DELETE FROM volumes
		WHERE container_id = $1
	`, id)
	return err
}
