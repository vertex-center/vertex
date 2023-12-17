package containersapi

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

func (c *Client) GetContainer(ctx context.Context, id types.ContainerID) (*types.Container, error) {
	var inst types.Container
	err := c.Request().
		Pathf("./container/%s", id).
		ToJSON(&inst).
		Fetch(ctx)
	return &inst, err
}

func (c *Client) DeleteContainer(ctx context.Context, id types.ContainerID) error {
	return c.Request().
		Pathf("./container/%s", id).
		Delete().
		Fetch(ctx)
}

func (c *Client) PatchContainer(ctx context.Context, id types.ContainerID, settings interface{}) error {
	return c.Request().
		Pathf("./container/%s", id).
		Patch().
		BodyJSON(&settings).
		Fetch(ctx)
}

func (c *Client) StartContainer(ctx context.Context, id types.ContainerID) error {
	return c.Request().
		Pathf("./container/%s/start", id).
		Post().
		Fetch(ctx)
}

func (c *Client) StopContainer(ctx context.Context, id types.ContainerID) error {
	return c.Request().
		Pathf("./container/%s/stop", id).
		Post().
		Fetch(ctx)
}

func (c *Client) AddContainerTag(ctx context.Context, id types.ContainerID, tagID types.TagID) error {
	return c.Request().
		Pathf("./container/%s/tag/%s", id, tagID).
		Post().
		Fetch(ctx)
}

func (c *Client) PatchContainerEnvironment(ctx context.Context, id types.ContainerID, env types.EnvVariables) error {
	return c.Request().
		Pathf("./container/%s/environment", id).
		Patch().
		BodyJSON(map[string]any{
			"env": env,
		}).
		Fetch(ctx)
}

func (c *Client) GetDocker(ctx context.Context, id types.ContainerID) (map[string]any, error) {
	var info map[string]any
	err := c.Request().
		Pathf("./container/%s/docker", id).
		ToJSON(&info).
		Fetch(ctx)
	return info, err
}

func (c *Client) RecreateDocker(ctx context.Context, id types.ContainerID) error {
	return c.Request().
		Pathf("./container/%s/docker/recreate", id).
		Post().
		Fetch(ctx)
}

func (c *Client) GetContainerLogs(ctx context.Context, id types.ContainerID) (string, error) {
	var logs string
	err := c.Request().
		Pathf("./container/%s/logs", id).
		ToJSON(&logs).
		Fetch(ctx)
	return logs, err
}

func (c *Client) UpdateServiceContainer(ctx context.Context, id types.ContainerID) error {
	return c.Request().
		Pathf("./container/%s/update/service", id).
		Post().
		Fetch(ctx)
}

func (c *Client) GetVersions(ctx context.Context, id types.ContainerID) ([]string, error) {
	var versions []string
	err := c.Request().
		Pathf("./container/%s/versions", id).
		ToJSON(&versions).
		Fetch(ctx)
	return versions, err
}

func (c *Client) WaitCondition(ctx context.Context, id types.ContainerID, condition container.WaitCondition) error {
	return c.Request().
		Pathf("./container/%s/wait/%s", id, condition).
		Fetch(ctx)
}
