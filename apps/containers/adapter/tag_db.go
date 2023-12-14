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

func (a *tagDBAdapter) GetTags(ctx context.Context) (types.Tags, error) {
	var tags types.Tags
	err := a.db.Select(&tags, `
		SELECT * FROM tags
	`)
	if errors.Is(err, sql.ErrNoRows) {
		return tags, nil
	}
	return tags, err
}

func (a *tagDBAdapter) CreateTag(ctx context.Context, tag types.Tag) error {
	_, err := a.db.NamedExec(`
		INSERT INTO tags (tag)
		VALUES (:tag)
	`, tag)
	return err
}

func (a *tagDBAdapter) DeleteTags(ctx context.Context, id types.TagID) error {
	_, err := a.db.Exec(`
		DELETE FROM tags
		WHERE id = $1
	`, id)
	return err
}
