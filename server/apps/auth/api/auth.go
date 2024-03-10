package api

import (
	"context"

	"github.com/vertex-center/vertex/server/apps/auth/core/types"
)

func (c *Client) Verify(ctx context.Context) (types.Session, error) {
	var session types.Session
	err := c.Request().
		Path("./verify").
		Post().
		ToJSON(&session).
		Fetch(ctx)
	return session, err
}
