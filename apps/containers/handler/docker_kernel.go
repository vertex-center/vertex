package handler

import (
	"bufio"
	"fmt"
	"io"

	"github.com/docker/docker/client"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
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

// docapi begin vx_containers_kernel_get_containers
// docapi method GET
// docapi summary Get containers
// docapi tags Apps/Containers
// docapi response 200 {[]Container} The containers.
// docapi response 500
// docapi end

func (h *DockerKernelHandler) GetContainers(c *router.Context) {
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

// docapi begin vx_containers_kernel_create_container
// docapi method POST
// docapi summary Create container
// docapi tags Apps/Containers
// docapi body CreateContainerOptions
// docapi response 200 {Container} The container.
// docapi response 500
// docapi end

func (h *DockerKernelHandler) CreateContainer(c *router.Context) {
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

// docapi begin vx_containers_kernel_delete_container
// docapi method DELETE
// docapi summary Delete container
// docapi tags Apps/Containers
// docapi query id {string} The container id.
// docapi response 200
// docapi response 404
// docapi response 500
// docapi end

func (h *DockerKernelHandler) DeleteContainer(c *router.Context) {
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

// docapi begin vx_containers_kernel_start_container
// docapi method POST
// docapi summary Start container
// docapi tags Apps/Containers
// docapi query id {string} The container id.
// docapi response 200
// docapi response 500
// docapi end

func (h *DockerKernelHandler) StartContainer(c *router.Context) {
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

// docapi begin vx_containers_kernel_stop_container
// docapi method POST
// docapi summary Stop container
// docapi tags Apps/Containers
// docapi query id {string} The container id.
// docapi response 200
// docapi response 500
// docapi end

func (h *DockerKernelHandler) StopContainer(c *router.Context) {
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

// docapi begin vx_containers_kernel_patch_container
// docapi method PATCH
// docapi summary Patch container
// docapi tags Apps/Containers
// docapi query id {string} The container id.
// docapi body PatchContainerOptions
// docapi response 200
// docapi response 500
// docapi end

func (h *DockerKernelHandler) InfoContainer(c *router.Context) {
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

// docapi begin vx_containers_kernel_logs_stdout_container
// docapi method GET
// docapi summary Get container stdout logs
// docapi desc Get container stdout logs as a stream.
// docapi tags Apps/Containers
// docapi query id {string} The container id.
// docapi response 200
// docapi response 500
// docapi end

func (h *DockerKernelHandler) LogsStdoutContainer(c *router.Context) {
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

// docapi begin vx_containers_kernel_logs_stderr_container
// docapi method GET
// docapi summary Get container stderr logs
// docapi desc Get container stderr logs as a stream.
// docapi tags Apps/Containers
// docapi query id {string} The container id.
// docapi response 200
// docapi response 500
// docapi end

func (h *DockerKernelHandler) LogsStderrContainer(c *router.Context) {
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

// docapi begin vx_containers_kernel_wait_container
// docapi method GET
// docapi summary Wait container
// docapi tags Apps/Containers
// docapi query id {string} The container id.
// docapi query cond {string} The condition to wait for.
// docapi response 200
// docapi response 500
// docapi end

func (h *DockerKernelHandler) WaitContainer(c *router.Context) {
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

// docapi begin vx_containers_kernel_info_image
// docapi method GET
// docapi summary Get image info
// docapi tags Apps/Containers
// docapi query id {string} The image id.
// docapi response 200 {InfoImageResponse} The image.
// docapi response 500
// docapi end

func (h *DockerKernelHandler) InfoImage(c *router.Context) {
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

// docapi begin vx_containers_kernel_pull_image
// docapi method POST
// docapi summary Pull image
// docapi desc Pull an image from a registry. The response is a stream of the logs.
// docapi tags Apps/Containers
// docapi body {PullImageOptions} The options to pull the image.
// docapi response 200
// docapi response 500
// docapi end

func (h *DockerKernelHandler) PullImage(c *router.Context) {
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

// docapi begin vx_containers_kernel_build_image
// docapi method POST
// docapi summary Build image
// docapi desc Build an image from a Dockerfile. The response is a stream of the logs.
// docapi tags Apps/Containers
// docapi body {BuildImageOptions} The options to build the image.
// docapi response 200
// docapi response 500
// docapi end

func (h *DockerKernelHandler) BuildImage(c *router.Context) {
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

// docapi begin vx_containers_kernel_create_volume
// docapi method POST
// docapi summary Create volume
// docapi tags Apps/Containers
// docapi body {CreateVolumeOptions} The options to create the volume.
// docapi response 200 {Volume} The volume.
// docapi response 500
// docapi end

func (h *DockerKernelHandler) CreateVolume(c *router.Context) {
	var options types.CreateVolumeOptions
	err := c.ParseBody(&options)
	if err != nil {
		return
	}

	vol, err := h.dockerService.CreateVolume(options)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToCreateVolume,
			PublicMessage:  fmt.Sprintf("Failed to create volume %s.", options.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(vol)
}

// docapi begin vx_containers_kernel_delete_volume
// docapi method DELETE
// docapi summary Delete volume
// docapi tags Apps/Containers
// docapi query name {string} The volume name.
// docapi response 200
// docapi response 500
// docapi end

func (h *DockerKernelHandler) DeleteVolume(c *router.Context) {
	name := c.Param("name")

	err := h.dockerService.DeleteVolume(name)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToDeleteVolume,
			PublicMessage:  fmt.Sprintf("Failed to delete volume %s.", name),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}