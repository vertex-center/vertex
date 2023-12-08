package api

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/apps/admin/handler"
	"github.com/vertex-center/vertex/pkg/user"
)

func (c *Client) GetSSHKeys(ctx context.Context) ([]types.PublicKey, error) {
	var keys []types.PublicKey
	err := c.Request().
		Pathf("./ssh").
		ToJSON(&keys).
		Fetch(ctx)
	return keys, err
}

func (c *KernelClient) GetSSHKeys(ctx context.Context) ([]types.PublicKey, error) {
	var keys []types.PublicKey
	err := c.Request().
		Pathf("./ssh").
		ToJSON(&keys).
		Fetch(ctx)
	return keys, err
}

func (c *Client) AddSSHKey(ctx context.Context, key string, username string) error {
	return c.Request().
		Pathf("./ssh").
		Post().
		BodyJSON(&handler.AddSSHKeyParams{
			AuthorizedKey: key,
			Username:      username,
		}).
		Fetch(ctx)
}

func (c *KernelClient) AddSSHKey(ctx context.Context, key string, username string) error {
	return c.Request().
		Pathf("./ssh").
		Post().
		BodyJSON(&handler.AddSSHKeyParams{
			AuthorizedKey: key,
			Username:      username,
		}).
		Fetch(ctx)
}

func (c *Client) DeleteSSHKey(ctx context.Context, key string, username string) error {
	return c.Request().
		Pathf("./ssh").
		BodyJSON(&handler.DeleteSSHKeyParams{
			Fingerprint: key,
			Username:    username,
		}).
		Delete().
		Fetch(ctx)
}

func (c *KernelClient) DeleteSSHKey(ctx context.Context, key string, username string) error {
	return c.Request().
		Pathf("./ssh").
		BodyJSON(&handler.DeleteSSHKeyParams{
			Fingerprint: key,
			Username:    username,
		}).
		Delete().
		Fetch(ctx)
}

func (c *Client) GetSSHUsers(ctx context.Context) ([]string, error) {
	var users []string
	err := c.Request().
		Pathf("./ssh/users").
		ToJSON(&users).
		Fetch(ctx)
	return users, err
}

func (c *KernelClient) GetSSHUsers(ctx context.Context) ([]user.User, error) {
	var users []user.User
	err := c.Request().
		Pathf("./ssh/users").
		ToJSON(&users).
		Fetch(ctx)
	return users, err
}
