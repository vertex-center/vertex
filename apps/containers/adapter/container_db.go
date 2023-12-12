package adapter

import (
	"context"

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

func (c *containerDBAdapter) GetContainer(ctx context.Context, id types.ContainerID) (*types.Container, error) {
	var container types.Container
	err := c.db.Get(&container, `
		SELECT * FROM containers
		WHERE id = $1
	`, id)
	return &container, err
}

func (c *containerDBAdapter) GetContainers(ctx context.Context) (types.Containers, error) {
	var containers types.Containers
	err := c.db.Select(&containers, `
		SELECT * FROM containers
	`)
	return containers, err
}

func (c *containerDBAdapter) CreateContainer(ctx context.Context, container types.Container) error {
	_, err := c.db.NamedExec(`
		INSERT INTO containers (id, service_id, user_id, image, image_tag, status, launch_on_startup, name, description, color, icon, command)
		VALUES (:id, :service_id, :user_id, :image, :image_tag, :status, :launch_on_startup, :name, :description, :color, :icon, :command)
	`, container)
	return err
}

func (c *containerDBAdapter) DeleteContainer(ctx context.Context, id types.ContainerID) error {
	_, err := c.db.Exec(`
		DELETE FROM containers
		WHERE id = $1
	`, id)
	return err
}

func (c *containerDBAdapter) GetTags(ctx context.Context) (types.Tags, error) {
	var tags types.Tags
	err := c.db.Select(&tags, `
		SELECT * FROM tags
	`)
	return tags, err
}
