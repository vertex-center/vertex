package containersapi

import (
	"context"
	"io"
	"net/http"

	"github.com/docker/docker/api/types/volume"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
)

func (c *KernelClient) CreateContainer(ctx context.Context, options types.CreateDockerContainerOptions) (types.CreateContainerResponse, error) {
	var res types.CreateContainerResponse
	err := c.Request().
		Pathf("./docker/containers").
		BodyJSON(options).
		ToJSON(&res).
		Fetch(ctx)
	return res, err
}

func (c *KernelClient) DeleteContainer(ctx context.Context, id string) error {
	return c.Request().
		Pathf("./docker/containers/%s", id).
		Delete().
		Fetch(ctx)
}

func (c *KernelClient) StartContainer(ctx context.Context, id string) error {
	return c.Request().
		Pathf("./docker/containers/%s/start", id).
		Post().
		Fetch(ctx)
}

func (c *KernelClient) StopContainer(ctx context.Context, id string) error {
	return c.Request().
		Pathf("./docker/containers/%s/stop", id).
		Post().
		Fetch(ctx)
}

func (c *KernelClient) GetContainerInfo(ctx context.Context, id string) (types.InfoContainerResponse, error) {
	var info types.InfoContainerResponse
	err := c.Request().
		Pathf("./docker/containers/%s/info", id).
		ToJSON(&info).
		Fetch(ctx)
	return info, err
}

func (c *KernelClient) GetImageInfo(ctx context.Context, id string) (types.InfoImageResponse, error) {
	var info types.InfoImageResponse
	err := c.Request().
		Pathf("./docker/images/%s/info", id).
		ToJSON(&info).
		Fetch(ctx)
	return info, err
}

func (c *KernelClient) BuildImage(ctx context.Context, options types.BuildImageOptions) (io.ReadCloser, error) {
	req, err := c.Request().
		Pathf("./docker/images/build").
		Post().
		BodyJSON(options).
		Request(ctx)

	if err != nil || req == nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (c *KernelClient) PullImage(ctx context.Context, options types.PullImageOptions) (io.ReadCloser, error) {
	req, err := c.Request().
		Pathf("./docker/images/pull").
		Post().
		BodyJSON(options).
		Request(ctx)

	if err != nil || req == nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (c *KernelClient) WaitContainer(ctx context.Context, id string, cond string) error {
	return c.Request().
		Pathf("./docker/containers/%s/wait/%s", id, cond).
		Fetch(ctx)
}

func (c *KernelClient) GetContainerStdout(ctx context.Context, id string) (io.ReadCloser, error) {
	var req *http.Request
	req, err := c.Request().
		Pathf("./docker/containers/%s/logs/stdout", id).
		Request(ctx)
	if err != nil || req == nil {
		return nil, err
	}

	stdout, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return stdout.Body, err
}

func (c *KernelClient) GetContainerStderr(ctx context.Context, id string) (io.ReadCloser, error) {
	var req *http.Request
	req, err := c.Request().
		Pathf("./docker/containers/%s/logs/stderr", id).
		Request(ctx)
	if err != nil || req == nil {
		return nil, err
	}

	stderr, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return stderr.Body, err
}

func (c *KernelClient) DeleteMounts(ctx context.Context, id string) error {
	return c.Request().
		Pathf("./docker/containers/%s/mounts", id).
		Delete().
		Fetch(ctx)
}

func (c *KernelClient) DeleteContainerVolumes(ctx context.Context, id string) error {
	return c.Request().
		Pathf("./docker/containers/%s/volumes", id).
		Delete().
		Fetch(ctx)
}

func (c *KernelClient) CreateVolume(ctx context.Context, name string) (volume.Volume, error) {
	var res volume.Volume
	err := c.Request().
		Path("./docker/volumes").
		BodyJSON(map[string]string{
			"name": name,
		}).
		ToJSON(&res).
		Post().
		Fetch(ctx)
	return res, err
}

func (c *KernelClient) DeleteVolume(ctx context.Context, name string) error {
	return c.Request().
		Path("./docker/volumes").
		BodyJSON(map[string]string{
			"name": name,
		}).
		Delete().
		Fetch(ctx)
}
