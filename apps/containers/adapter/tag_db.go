package adapter

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/storage"
)

type tagDBAdapter struct {
	db storage.DB
}

func NewTagDBAdapter(db storage.DB) port.TagAdapter {
	return &tagDBAdapter{db}
}

func (a *tagDBAdapter) CreateTags(ctx context.Context, tags types.Tags) error {
	for _, t := range tags {
		_, err := a.db.NamedExec(`
			INSERT INTO tags (container_id, tag)
			VALUES (:container_id, :tag)
		`, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *tagDBAdapter) DeleteTags(ctx context.Context, id types.ContainerID) error {
	_, err := a.db.Exec(`
		DELETE FROM tags
		WHERE container_id = $1
	`, id)
	return err
}

func (a *tagDBAdapter) GetContainerTags(ctx context.Context, id types.ContainerID) (types.Tags, error) {
	var tags types.Tags
	err := a.db.Select(&tags, `
		SELECT * FROM tags
		WHERE container_id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return tags, nil
	}
	return tags, err
}

func (a *tagDBAdapter) GetUniqueTags(ctx context.Context) (types.Tags, error) {
	var tags types.Tags
	err := a.db.Select(&tags, `
		SELECT DISTINCT tag FROM tags
	`)
	if errors.Is(err, sql.ErrNoRows) {
		return tags, nil
	}
	return tags, err
}