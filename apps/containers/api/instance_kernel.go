package containersapi

import (
	"context"
	"io"
	"net/http"

	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types/api"
)

func (c *KernelClient) CreateContainer(ctx context.Context, options types.CreateContainerOptions) (types.CreateContainerResponse, *api.Error) {
	var apiError api.Error
	var res types.CreateContainerResponse
	err := c.Request().
		Pathf("./docker/container").
		BodyJSON(options).
		ErrorJSON(&apiError).
		ToJSON(&res).
		Fetch(ctx)
	return res, api.HandleError(err, apiError)
}

func (c *KernelClient) DeleteContainer(ctx context.Context, id string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./docker/container/%s", id).
		Delete().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *KernelClient) StartContainer(ctx context.Context, id string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./docker/container/%s/start", id).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *KernelClient) StopContainer(ctx context.Context, id string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./docker/container/%s/stop", id).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *KernelClient) GetContainerInfo(ctx context.Context, id string) (types.InfoContainerResponse, *api.Error) {
	var apiError api.Error
	var info types.InfoContainerResponse
	err := c.Request().
		Pathf("./docker/container/%s/info", id).
		ErrorJSON(&apiError).
		ToJSON(&info).
		Fetch(ctx)
	return info, api.HandleError(err, apiError)
}

func (c *KernelClient) GetImageInfo(ctx context.Context, id string) (types.InfoImageResponse, *api.Error) {
	var apiError api.Error
	var info types.InfoImageResponse
	err := c.Request().
		Pathf("./docker/image/%s/info", id).
		ErrorJSON(&apiError).
		ToJSON(&info).
		Fetch(ctx)
	return info, api.HandleError(err, apiError)
}

func (c *KernelClient) BuildImage(ctx context.Context, options types.BuildImageOptions) (io.ReadCloser, *api.Error) {
	var apiError api.Error
	req, err := c.Request().
		Pathf("./docker/image/build").
		Post().
		BodyJSON(options).
		ErrorJSON(&apiError).
		Request(ctx)

	if err != nil || req == nil {
		return nil, api.HandleError(err, apiError)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, api.HandleError(err, apiError)
	}
	return res.Body, nil
}

func (c *KernelClient) PullImage(ctx context.Context, options types.PullImageOptions) (io.ReadCloser, *api.Error) {
	var apiError api.Error
	req, err := c.Request().
		Pathf("./docker/image/pull").
		Post().
		BodyJSON(options).
		ErrorJSON(&apiError).
		Request(ctx)

	if err != nil || req == nil {
		return nil, api.HandleError(err, apiError)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, api.HandleError(err, apiError)
	}
	return res.Body, nil
}

func (c *KernelClient) WaitContainer(ctx context.Context, id string, cond string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./docker/container/%s/wait/%s", id, cond).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *KernelClient) GetContainerStdout(ctx context.Context, id string) (io.ReadCloser, *api.Error) {
	var apiError api.Error
	var req *http.Request
	req, err := c.Request().
		Pathf("./docker/container/%s/logs/stdout", id).
		ErrorJSON(&apiError).
		Request(ctx)
	if err != nil || req == nil {
		return nil, api.HandleError(err, apiError)
	}

	stdout, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, api.HandleError(err, apiError)
	}

	return stdout.Body, api.HandleError(err, apiError)
}

func (c *KernelClient) GetContainerStderr(ctx context.Context, id string) (io.ReadCloser, *api.Error) {
	var apiError api.Error
	var req *http.Request
	req, err := c.Request().
		Pathf("./docker/container/%s/logs/stderr", id).
		ErrorJSON(&apiError).
		Request(ctx)
	if err != nil || req == nil {
		return nil, api.HandleError(err, apiError)
	}

	stderr, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, api.HandleError(err, apiError)
	}

	return stderr.Body, api.HandleError(err, apiError)
}

func (c *KernelClient) DeleteMounts(ctx context.Context, id string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./docker/container/%s/mounts", id).
		Delete().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}
