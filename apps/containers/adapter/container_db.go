package adapter

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/storage"
)

type containerDBAdapter struct {
	db storage.DB
}

func NewContainerDBAdapter(db storage.DB) port.ContainerAdapter {
	return &containerDBAdapter{db}
}

func (a *containerDBAdapter) GetContainer(ctx context.Context, id types.ContainerID) (*types.Container, error) {
	var container types.Container
	err := a.db.Get(&container, `
		SELECT * FROM containers
		WHERE id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return &container, nil
	}
	return &container, err
}

func (a *containerDBAdapter) GetContainers(ctx context.Context) (types.Containers, error) {
	var containers types.Containers
	err := a.db.Select(&containers, `
		SELECT * FROM containers
	`)
	if errors.Is(err, sql.ErrNoRows) {
		return containers, nil
	}
	return containers, err
}

func (a *containerDBAdapter) CreateContainer(ctx context.Context, c types.Container) error {
	_, err := a.db.NamedExec(`
		INSERT INTO containers (id, service_id, user_id, image, image_tag, status, launch_on_startup, name, description, color, icon, command)
		VALUES (:id, :service_id, :user_id, :image, :image_tag, :status, :launch_on_startup, :name, :description, :color, :icon, :command)
	`, c)
	return err
}

func (a *containerDBAdapter) DeleteContainer(ctx context.Context, id types.ContainerID) error {
	_, err := a.db.Exec(`
		DELETE FROM containers
		WHERE id = $1
	`, id)
	return err
}
