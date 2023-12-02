package api

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/core/service"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/apps/admin/handler"
	"github.com/vertex-center/vertex/core/types/api"
)

func (c *Client) GetSSHKeys(ctx context.Context) ([]types.PublicKey, error) {
	var keys []types.PublicKey
	var apiError api.Error

	err := c.Request().
		Pathf("./ssh").
		ToJSON(&keys).
		ErrorJSON(&apiError).
		Fetch(ctx)

	return keys, err
}

func (c *Client) AddSSHKey(ctx context.Context, key string, username string) error {
	var apiError api.Error
	err := c.Request().
		Pathf("./ssh").
		Post().
		BodyJSON(&handler.AddSSHKeyBody{
			AuthorizedKey: key,
			Username:      username,
		}).
		ErrorJSON(&apiError).
		Fetch(ctx)

	if apiError.Code == api.ErrFailedToAddSSHKey {
		return service.ErrFailedToAddKey
	}
	if apiError.Code == api.ErrUserNotFound {
		return service.ErrUserNotFound
	}
	return err
}

func (c *Client) DeleteSSHKey(ctx context.Context, key string, username string) error {
	var apiError api.Error
	err := c.Request().
		Pathf("./ssh").
		BodyJSON(&handler.DeleteSSHKeyBody{
			Fingerprint: key,
			Username:    username,
		}).
		Delete().
		ErrorJSON(&apiError).
		Fetch(ctx)

	if apiError.Code == api.ErrUserNotFound {
		return service.ErrUserNotFound
	}
	return err
}

func (c *Client) GetSSHUsers(ctx context.Context) ([]string, error) {
	var users []string
	var apiError api.Error

	err := c.Request().
		Pathf("./ssh/users").
		ToJSON(&users).
		ErrorJSON(&apiError).
		Fetch(ctx)

	return users, err
}
