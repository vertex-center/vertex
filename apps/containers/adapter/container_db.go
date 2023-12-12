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
	tx, err := a.db.Beginx()
	if err != nil {
		return nil, err
	}

	var container types.Container
	err = tx.Get(&container, `
		SELECT * FROM containers
		WHERE id = $1
	`, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Select(&container.Env, `
		SELECT * FROM env_variables
		WHERE container_id = $1
	`, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Select(&container.Capabilities, `
		SELECT * FROM capabilities
		WHERE container_id = $1
	`, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Select(&container.Ports, `
		SELECT * FROM ports
		WHERE container_id = $1
	`, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Select(&container.Volumes, `
		SELECT * FROM volumes
		WHERE container_id = $1
	`, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Select(&container.Sysctls, `
		SELECT * FROM sysctls
		WHERE container_id = $1
	`, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Select(&container.Tags, `
		SELECT * FROM tags
		WHERE container_id = $1
	`, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
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
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		INSERT INTO containers (id, service_id, user_id, image, image_tag, status, launch_on_startup, name, description, color, icon, command)
		VALUES (:id, :service_id, :user_id, :image, :image_tag, :status, :launch_on_startup, :name, :description, :color, :icon, :command)
	`, c)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, e := range c.Env {
		_, err = tx.NamedExec(`
			INSERT INTO env_variables (container_id, type, name, value, default_value, description)
			VALUES (:container_id, :type, :name, :value, :default_value, :description)
		`, e)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	for _, cp := range c.Capabilities {
		_, err = tx.NamedExec(`
			INSERT INTO capabilities (container_id, name)
			VALUES (:container_id, :name)
		`, cp)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	for _, p := range c.Ports {
		_, err = tx.NamedExec(`
			INSERT INTO ports (container_id, internal_port, external_port)
			VALUES (:container_id, :internal_port, :external_port)
		`, p)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	for _, v := range c.Volumes {
		_, err = tx.NamedExec(`
			INSERT INTO volumes (container_id, internal_path, external_path)
			VALUES (:container_id, :internal_path, :external_path)
		`, v)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	for _, s := range c.Sysctls {
		_, err = tx.NamedExec(`
			INSERT INTO sysctls (container_id, name, value)
			VALUES (:container_id, :name, :value)
		`, s)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	for _, t := range c.Tags {
		_, err = tx.NamedExec(`
			INSERT INTO tags (container_id, tag)
			VALUES (:container_id, :tag)
		`, t)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (a *containerDBAdapter) DeleteContainer(ctx context.Context, id types.ContainerID) error {
	_, err := a.db.Exec(`
		DELETE FROM containers
		WHERE id = $1
	`, id)
	return err
}

func (a *containerDBAdapter) GetTags(ctx context.Context) (types.Tags, error) {
	var tags types.Tags
	err := a.db.Select(&tags, `
		SELECT * FROM tags
	`)
	return tags, err
}
