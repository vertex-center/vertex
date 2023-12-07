package handler

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type dockerKernelHandler struct {
	dockerService port.DockerService
}

func NewDockerKernelHandler(dockerKernelService port.DockerService) port.DockerKernelHandler {
	return &dockerKernelHandler{
		dockerService: dockerKernelService,
	}
}

func (h *dockerKernelHandler) GetContainers(c *router.Context) {
	containers, err := h.dockerService.ListContainers()
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToListContainers,
			PublicMessage:  "Failed to list containers.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(containers)
}

func (h *dockerKernelHandler) GetContainersInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get containers"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel([]types.Container{}),
		),
	}
}

func (h *dockerKernelHandler) CreateContainer(c *router.Context) {
	var options types.CreateContainerOptions
	err := c.ParseBody(&options)
	if err != nil {
		return
	}

	res, err := h.dockerService.CreateContainer(options)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToCreateContainer,
			PublicMessage:  "Failed to create container.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(res)
}

func (h *dockerKernelHandler) CreateContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Create container"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.Container{}),
		),
	}
}

func (h *dockerKernelHandler) DeleteContainer(c *router.Context) {
	id := c.Param("id")

	err := h.dockerService.DeleteContainer(id)
	if err != nil && client.IsErrNotFound(err) {
		c.NotFound(router.Error{
			Code:           types.ErrCodeContainerNotFound,
			PublicMessage:  fmt.Sprintf("Container %s not found.", id),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToDeleteContainer,
			PublicMessage:  fmt.Sprintf("Failed to delete container %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *dockerKernelHandler) DeleteContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Delete container"),
		oapi.Response(http.StatusNoContent),
	}
}

func (h *dockerKernelHandler) StartContainer(c *router.Context) {
	id := c.Param("id")

	err := h.dockerService.StartContainer(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToStartContainer,
			PublicMessage:  fmt.Sprintf("Failed to start container %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *dockerKernelHandler) StartContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Start container"),
		oapi.Response(http.StatusNoContent),
	}
}

func (h *dockerKernelHandler) StopContainer(c *router.Context) {
	id := c.Param("id")

	err := h.dockerService.StopContainer(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToStopContainer,
			PublicMessage:  fmt.Sprintf("Failed to stop container %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *dockerKernelHandler) StopContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Stop container"),
		oapi.Response(http.StatusNoContent),
	}
}

func (h *dockerKernelHandler) InfoContainer(c *router.Context) {
	id := c.Param("id")

	info, err := h.dockerService.InfoContainer(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetContainerInfo,
			PublicMessage:  fmt.Sprintf("Failed to get info for container %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(info)
}

func (h *dockerKernelHandler) InfoContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get container info"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.Container{}),
		),
	}
}

func (h *dockerKernelHandler) LogsStdoutContainer(c *router.Context) {
	id := c.Param("id")

	stdout, err := h.dockerService.LogsStdoutContainer(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetContainerLogs,
			PublicMessage:  fmt.Sprintf("Failed to get logs for container %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}
	defer stdout.Close()

	scanner := bufio.NewScanner(stdout)

	c.Stream(func(w io.Writer) bool {
		if scanner.Err() != nil {
			return false
		}
		if !scanner.Scan() {
			return false
		}

		_, err := fmt.Fprintln(w, scanner.Text())
		if err != nil {
			log.Error(err)
			return false
		}
		return true
	})
}

func (h *dockerKernelHandler) LogsStdoutContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get container stdout logs"),
		oapi.Description("Get container stdout logs as a stream."),
		oapi.Response(http.StatusOK),
	}
}

func (h *dockerKernelHandler) LogsStderrContainer(c *router.Context) {
	id := c.Param("id")

	stderr, err := h.dockerService.LogsStderrContainer(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetContainerLogs,
			PublicMessage:  fmt.Sprintf("Failed to get logs for container %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}
	defer stderr.Close()

	scanner := bufio.NewScanner(stderr)

	c.Stream(func(w io.Writer) bool {
		if scanner.Err() != nil {
			return false
		}
		if !scanner.Scan() {
			return false
		}

		_, err := fmt.Fprintln(w, scanner.Text())
		if err != nil {
			log.Error(err)
			return false
		}
		return true
	})
}

func (h *dockerKernelHandler) LogsStderrContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get container stderr logs"),
		oapi.Description("Get container stderr logs as a stream."),
		oapi.Response(http.StatusOK),
	}
}

func (h *dockerKernelHandler) WaitContainer(c *router.Context) {
	id := c.Param("id")
	cond := c.Param("cond")

	err := h.dockerService.WaitContainer(id, types.WaitContainerCondition(cond))
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToWaitContainer,
			PublicMessage:  fmt.Sprintf("Failed to wait the event '%s' for container %s.", cond, id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *dockerKernelHandler) WaitContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Wait container"),
		oapi.Response(http.StatusNoContent),
	}
}

func (h *dockerKernelHandler) DeleteMounts(c *router.Context) {
	id := c.Param("id")

	err := h.dockerService.DeleteMounts(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToDeleteMounts,
			PublicMessage:  fmt.Sprintf("Failed to delete mounts of %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *dockerKernelHandler) DeleteMountsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Delete mounts"),
		oapi.Response(http.StatusNoContent),
	}
}

func (h *dockerKernelHandler) InfoImage(c *router.Context) {
	id := c.Param("id")

	info, err := h.dockerService.InfoImage(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetImageInfo,
			PublicMessage:  fmt.Sprintf("Failed to get info for image %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(info)
}

func (h *dockerKernelHandler) InfoImageInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get image info"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.InfoImageResponse{}),
		),
	}
}

func (h *dockerKernelHandler) PullImage(c *router.Context) {
	var options types.PullImageOptions
	err := c.ParseBody(&options)
	if err != nil {
		return
	}

	r, err := h.dockerService.PullImage(options)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToPullImage,
			PublicMessage:  "Failed to pull image.",
			PrivateMessage: err.Error(),
		})
		return
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)

	c.Stream(func(w io.Writer) bool {
		if scanner.Err() != nil {
			return false
		}
		if !scanner.Scan() {
			return false
		}

		_, err := fmt.Fprintln(w, scanner.Text())
		if err != nil {
			log.Error(err)
			return false
		}
		return true
	})
}

func (h *dockerKernelHandler) PullImageInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Pull image"),
		oapi.Response(http.StatusOK),
	}
}

func (h *dockerKernelHandler) BuildImage(c *router.Context) {
	var options types.BuildImageOptions
	err := c.ParseBody(&options)
	if err != nil {
		return
	}

	res, err := h.dockerService.BuildImage(options)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToBuildImage,
			PublicMessage:  "Failed to build image.",
			PrivateMessage: err.Error(),
		})
		return
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)

	c.Stream(func(w io.Writer) bool {
		if scanner.Err() != nil {
			log.Error(scanner.Err())
			return false
		}

		if !scanner.Scan() {
			return false
		}

		_, err := io.WriteString(w, scanner.Text()+"\n")
		if err != nil {
			log.Error(err)
			return false
		}
		return true
	})
}

func (h *dockerKernelHandler) BuildImageInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Build image"),
		oapi.Response(http.StatusOK),
	}
}
