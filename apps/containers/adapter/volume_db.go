package adapter

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/storage"
)

type volumeDBAdapter struct {
	db storage.DB
}

func NewVolumeDBAdapter(db storage.DB) port.VolumeAdapter {
	return &volumeDBAdapter{db}
}

func (a *volumeDBAdapter) GetVolumes(ctx context.Context, id types.ContainerID) (types.Volumes, error) {
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

func (a *volumeDBAdapter) CreateVolumes(ctx context.Context, volumes types.Volumes) error {
	for _, v := range volumes {
		_, err := a.db.NamedExec(`
			INSERT INTO volumes (container_id, internal_path, external_path)
			VALUES (:container_id, :internal_path, :external_path)
		`, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *volumeDBAdapter) DeleteVolumes(ctx context.Context, id types.ContainerID) error {
	_, err := a.db.Exec(`
		DELETE FROM volumes
		WHERE container_id = $1
	`, id)
	return err
}
