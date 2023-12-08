package containersapi

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

func (c *Client) GetContainer(ctx context.Context, uuid uuid.UUID) (*types.Container, error) {
	var inst types.Container
	err := c.Request().
		Pathf("./container/%s", uuid).
		ToJSON(&inst).
		Fetch(ctx)
	return &inst, err
}

func (c *Client) DeleteContainer(ctx context.Context, uuid uuid.UUID) error {
	return c.Request().
		Pathf("./container/%s", uuid).
		Delete().
		Fetch(ctx)
}

func (c *Client) PatchContainer(ctx context.Context, uuid uuid.UUID, settings types.ContainerSettings) error {
	return c.Request().
		Pathf("./container/%s", uuid).
		Patch().
		BodyJSON(&settings).
		Fetch(ctx)
}

func (c *Client) StartContainer(ctx context.Context, uuid uuid.UUID) error {
	return c.Request().
		Pathf("./container/%s/start", uuid).
		Post().
		Fetch(ctx)
}

func (c *Client) StopContainer(ctx context.Context, uuid uuid.UUID) error {
	return c.Request().
		Pathf("./container/%s/stop", uuid).
		Post().
		Fetch(ctx)
}

func (c *Client) PatchContainerEnvironment(ctx context.Context, uuid uuid.UUID, env map[string]string) error {
	return c.Request().
		Pathf("./container/%s/environment", uuid).
		Patch().
		BodyJSON(map[string]any{
			"env": env,
		}).
		Fetch(ctx)
}

func (c *Client) GetDocker(ctx context.Context, uuid uuid.UUID) (map[string]any, error) {
	var info map[string]any
	err := c.Request().
		Pathf("./container/%s/docker", uuid).
		ToJSON(&info).
		Fetch(ctx)
	return info, err
}

func (c *Client) RecreateDocker(ctx context.Context, uuid uuid.UUID) error {
	return c.Request().
		Pathf("./container/%s/docker/recreate", uuid).
		Post().
		Fetch(ctx)
}

func (c *Client) GetContainerLogs(ctx context.Context, uuid uuid.UUID) (string, error) {
	var logs string
	err := c.Request().
		Pathf("./container/%s/logs", uuid).
		ToJSON(&logs).
		Fetch(ctx)
	return logs, err
}

func (c *Client) UpdateServiceContainer(ctx context.Context, uuid uuid.UUID) error {
	return c.Request().
		Pathf("./container/%s/update/service", uuid).
		Post().
		Fetch(ctx)
}

func (c *Client) GetVersions(ctx context.Context, uuid uuid.UUID) ([]string, error) {
	var versions []string
	err := c.Request().
		Pathf("./container/%s/versions", uuid).
		ToJSON(&versions).
		Fetch(ctx)
	return versions, err
}

func (c *Client) WaitCondition(ctx context.Context, uuid uuid.UUID, condition container.WaitCondition) error {
	return c.Request().
		Pathf("./container/%s/wait/%s", uuid, condition).
		Fetch(ctx)
}
