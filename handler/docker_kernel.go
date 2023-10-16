package handler

import (
	"bufio"
	"fmt"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"io"

	"github.com/docker/docker/client"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
)

type DockerKernelHandler struct {
	dockerService port.DockerService
}

func NewDockerKernelHandler(dockerKernelService port.DockerService) port.DockerKernelHandler {
	return &DockerKernelHandler{
		dockerService: dockerKernelService,
	}
}

func (h *DockerKernelHandler) GetContainers(c *router.Context) {
	containers, err := h.dockerService.ListContainers()
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToListContainers,
			PublicMessage:  "Failed to list containers.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(containers)
}

func (h *DockerKernelHandler) CreateContainer(c *router.Context) {
	var options types.CreateContainerOptions
	err := c.ParseBody(&options)
	if err != nil {
		return
	}

	res, err := h.dockerService.CreateContainer(options)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToCreateContainer,
			PublicMessage:  "Failed to create container.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(res)
}

func (h *DockerKernelHandler) DeleteContainer(c *router.Context) {
	id := c.Param("id")

	err := h.dockerService.DeleteContainer(id)
	if err != nil && client.IsErrNotFound(err) {
		c.NotFound(router.Error{
			Code:           api.ErrContainerNotFound,
			PublicMessage:  fmt.Sprintf("Container %s not found.", id),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToDeleteContainer,
			PublicMessage:  fmt.Sprintf("Failed to delete container %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *DockerKernelHandler) StartContainer(c *router.Context) {
	id := c.Param("id")

	err := h.dockerService.StartContainer(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToStartContainer,
			PublicMessage:  fmt.Sprintf("Failed to start container %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *DockerKernelHandler) StopContainer(c *router.Context) {
	id := c.Param("id")

	err := h.dockerService.StopContainer(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToStopContainer,
			PublicMessage:  fmt.Sprintf("Failed to stop container %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *DockerKernelHandler) InfoContainer(c *router.Context) {
	id := c.Param("id")

	info, err := h.dockerService.InfoContainer(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetContainerInfo,
			PublicMessage:  fmt.Sprintf("Failed to get info for container %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(info)
}

func (h *DockerKernelHandler) LogsStdoutContainer(c *router.Context) {
	id := c.Param("id")

	stdout, err := h.dockerService.LogsStdoutContainer(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetContainerLogs,
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

func (h *DockerKernelHandler) LogsStderrContainer(c *router.Context) {
	id := c.Param("id")

	stderr, err := h.dockerService.LogsStderrContainer(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetContainerLogs,
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

func (h *DockerKernelHandler) WaitContainer(c *router.Context) {
	id := c.Param("id")
	cond := c.Param("cond")

	err := h.dockerService.WaitContainer(id, types.WaitContainerCondition(cond))
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToWaitContainer,
			PublicMessage:  fmt.Sprintf("Failed to wait the event '%s' for container %s.", cond, id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *DockerKernelHandler) InfoImage(c *router.Context) {
	id := c.Param("id")

	info, err := h.dockerService.InfoImage(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetImageInfo,
			PublicMessage:  fmt.Sprintf("Failed to get info for image %s.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(info)
}

func (h *DockerKernelHandler) PullImage(c *router.Context) {
	var options types.PullImageOptions
	err := c.ParseBody(&options)
	if err != nil {
		return
	}

	r, err := h.dockerService.PullImage(options)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToPullImage,
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

func (h *DockerKernelHandler) BuildImage(c *router.Context) {
	var options types.BuildImageOptions
	err := c.ParseBody(&options)
	if err != nil {
		return
	}

	res, err := h.dockerService.BuildImage(options)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToBuildImage,
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
