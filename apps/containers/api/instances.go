package containersapi

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

func (c *Client) GetContainer(ctx context.Context, id uuid.UUID) (*types.Container, error) {
	var inst types.Container
	err := c.Request().
		Pathf("./containers/%s", id).
		ToJSON(&inst).
		Fetch(ctx)
	return &inst, err
}

func (c *Client) GetContainers(ctx context.Context) (map[uuid.UUID]*types.Container, error) {
	var insts map[uuid.UUID]*types.Container
	err := c.Request().
		Path("./containers").
		ToJSON(&insts).
		Fetch(ctx)
	return insts, err
}

func (c *Client) CreateContainer(ctx context.Context, templateID string) (types.Container, error) {
	var ctr types.Container
	err := c.Request().
		Path("./containers").
		BodyJSON(map[string]string{
			"template_id": templateID,
		}).
		Post().
		ToJSON(&ctr).
		Fetch(ctx)
	return ctr, err
}

func (c *Client) DeleteContainer(ctx context.Context, id uuid.UUID) error {
	return c.Request().
		Pathf("./containers/%s", id).
		Delete().
		Fetch(ctx)
}

func (c *Client) PatchContainer(ctx context.Context, id uuid.UUID, settings interface{}) error {
	return c.Request().
		Pathf("./containers/%s", id).
		Patch().
		BodyJSON(&settings).
		Fetch(ctx)
}

func (c *Client) StartContainer(ctx context.Context, id uuid.UUID) error {
	return c.Request().
		Pathf("./containers/%s/start", id).
		Post().
		Fetch(ctx)
}

func (c *Client) StopContainer(ctx context.Context, id uuid.UUID) error {
	return c.Request().
		Pathf("./containers/%s/stop", id).
		Post().
		Fetch(ctx)
}

func (c *Client) AddContainerTag(ctx context.Context, id uuid.UUID, tagID uuid.UUID) error {
	return c.Request().
		Pathf("./containers/%s/tags/%s", id, tagID).
		Put().
		Fetch(ctx)
}

func (c *Client) PatchContainerEnvironment(ctx context.Context, id uuid.UUID, env types.EnvVariables) error {
	return c.Request().
		Pathf("./containers/%s/environment", id).
		Patch().
		BodyJSON(map[string]any{
			"env": env,
		}).
		Fetch(ctx)
}

func (c *Client) GetDocker(ctx context.Context, id uuid.UUID) (map[string]any, error) {
	var info map[string]any
	err := c.Request().
		Pathf("./containers/%s/docker", id).
		ToJSON(&info).
		Fetch(ctx)
	return info, err
}

func (c *Client) RecreateDocker(ctx context.Context, id uuid.UUID) error {
	return c.Request().
		Pathf("./containers/%s/docker/recreate", id).
		Post().
		Fetch(ctx)
}

func (c *Client) GetContainerLogs(ctx context.Context, id uuid.UUID) (string, error) {
	var logs string
	err := c.Request().
		Pathf("./containers/%s/logs", id).
		ToJSON(&logs).
		Fetch(ctx)
	return logs, err
}

func (c *Client) UpdateServiceContainer(ctx context.Context, id uuid.UUID) error {
	return c.Request().
		Pathf("./containers/%s/update/service", id).
		Post().
		Fetch(ctx)
}

func (c *Client) GetVersions(ctx context.Context, id uuid.UUID) ([]string, error) {
	var versions []string
	err := c.Request().
		Pathf("./containers/%s/versions", id).
		ToJSON(&versions).
		Fetch(ctx)
	return versions, err
}

func (c *Client) WaitCondition(ctx context.Context, id uuid.UUID, condition container.WaitCondition) error {
	return c.Request().
		Pathf("./containers/%s/wait/%s", id, condition).
		Fetch(ctx)
}

func (c *Client) CheckForUpdates(ctx context.Context) ([]types.Container, error) {
	var insts []types.Container
	err := c.Request().
		Path("./containers/checkupdates").
		ToJSON(&insts).
		Fetch(ctx)
	return insts, err
}
