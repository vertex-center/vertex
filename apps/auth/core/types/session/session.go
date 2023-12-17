package session

import (
	"context"

	"github.com/vertex-center/vertex/common/uuid"
)

type Session struct {
	UserID uuid.UUID
}

func Get(ctx context.Context) Session {
	return Session{
		UserID: ctx.Value("user_id").(uuid.UUID),
	}
}
