package adapter

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vertex/common/uuid"
)

type tagDBAdapter struct {
	db storage.DB
}

func NewTagDBAdapter(db storage.DB) port.TagAdapter {
	return &tagDBAdapter{db}
}

func (a *tagDBAdapter) GetTag(ctx context.Context, userID uuid.UUID, name string) (types.Tag, error) {
	var tag types.Tag
	err := a.db.Get(&tag, `
		SELECT * FROM tags
		WHERE user_id = $1 AND name = $2
	`, userID, name)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.NotFoundf("tag %s", name)
	}
	return tag, err
}

func (a *tagDBAdapter) GetTags(ctx context.Context, userID uuid.UUID) (types.Tags, error) {
	var tags types.Tags
	err := a.db.Select(&tags, `
		SELECT * FROM tags
		WHERE user_id = $1
	`, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return tags, nil
	}
	return tags, err
}

func (a *tagDBAdapter) CreateTag(ctx context.Context, tag types.Tag) error {
	_, err := a.db.NamedExec(`
		INSERT INTO tags (id, user_id, name)
		VALUES (:id, :user_id, :name)
	`, tag)
	return err
}

func (a *tagDBAdapter) DeleteTag(ctx context.Context, id types.TagID) error {
	_, err := a.db.Exec(`
		DELETE FROM tags
		WHERE id = $1
	`, id)
	return err
}
