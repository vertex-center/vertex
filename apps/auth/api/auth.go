package api

import (
	"context"

	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/core/types/api"
)

func (c *Client) Verify(ctx context.Context) (types.Session, *api.Error) {
	var apiError api.Error
	var session types.Session
	err := c.Request().
		Path("./verify").
		Post().
		ErrorJSON(&apiError).
		ToJSON(&session).
		Fetch(ctx)
	return session, api.HandleError(err, apiError)
}
