package api

import (
	"context"

	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/handler"
	"github.com/vertex-center/vertex/pkg/user"
)

func (c *Client) GetSSHKeys(ctx context.Context) ([]types.PublicKey, error) {
	var keys []types.PublicKey
	var apiError api.Error

	err := c.Request().
		Pathf("./security/ssh").
		ToJSON(&keys).
		ErrorJSON(&apiError).
		Fetch(ctx)

	return keys, err
}

func (c *Client) AddSSHKey(ctx context.Context, key string) error {
	var apiError api.Error
	err := c.Request().
		Pathf("./security/ssh").
		Post().
		BodyJSON(&handler.AddSSHKeyBody{
			AuthorizedKey: key,
		}).
		ErrorJSON(&apiError).
		Fetch(ctx)

	return err
}

func (c *Client) DeleteSSHKey(ctx context.Context, key string) error {
	var apiError api.Error
	err := c.Request().
		Pathf("./security/ssh/%s", key).
		Delete().
		ErrorJSON(&apiError).
		Fetch(ctx)

	return err
}

func (c *Client) GetSSHUsers(ctx context.Context) ([]user.User, error) {
	var users []user.User
	var apiError api.Error

	err := c.Request().
		Pathf("./security/ssh/users").
		ToJSON(&users).
		ErrorJSON(&apiError).
		Fetch(ctx)

	return users, err
}
