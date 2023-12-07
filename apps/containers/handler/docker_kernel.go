package handler

import (
	"bufio"
	"fmt"
	"io"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	apierrors "github.com/juju/errors"
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

func (h *dockerKernelHandler) GetContainers() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]types.DockerContainer, error) {
		return h.dockerService.ListContainers()
	})
}

func (h *dockerKernelHandler) GetContainersInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getContainers"),
		oapi.Summary("Get containers"),
	}
}

func (h *dockerKernelHandler) CreateContainer() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *types.CreateContainerOptions) (*types.CreateContainerResponse, error) {
		res, err := h.dockerService.CreateContainer(*params)
		return &res, err
	})
}

func (h *dockerKernelHandler) CreateContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("createContainer"),
		oapi.Summary("Create container"),
	}
}

type DeleteDockerContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) DeleteContainer() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *DeleteDockerContainerParams) error {
		err := h.dockerService.DeleteContainer(params.ID)
		if err != nil && client.IsErrNotFound(err) {
			return apierrors.NewNotFound(err, "container not found")
		}
		return err
	})
}

func (h *dockerKernelHandler) DeleteContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("deleteContainer"),
		oapi.Summary("Delete container"),
	}
}

type StartDockerContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) StartContainer() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *StartDockerContainerParams) error {
		err := h.dockerService.StartContainer(params.ID)
		if err != nil && client.IsErrNotFound(err) {
			return apierrors.NewNotFound(err, "container not found")
		}
		return err
	})
}

func (h *dockerKernelHandler) StartContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("startContainer"),
		oapi.Summary("Start container"),
	}
}

type StopDockerContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) StopContainer() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *StopDockerContainerParams) error {
		err := h.dockerService.StopContainer(params.ID)
		if err != nil && client.IsErrNotFound(err) {
			return apierrors.NewNotFound(err, "container not found")
		}
		return err
	})
}

func (h *dockerKernelHandler) StopContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("stopContainer"),
		oapi.Summary("Stop container"),
	}
}

type InfoContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) InfoContainer() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *InfoContainerParams) (*types.InfoContainerResponse, error) {
		info, err := h.dockerService.InfoContainer(params.ID)
		if err != nil && client.IsErrNotFound(err) {
			return nil, apierrors.NewNotFound(err, "container not found")
		}
		return &info, err
	})
}

func (h *dockerKernelHandler) InfoContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("infoContainer"),
		oapi.Summary("Get container info"),
	}
}

type LogsStdoutContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) LogsStdoutContainer() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *LogsStdoutContainerParams) error {
		stdout, err := h.dockerService.LogsStdoutContainer(params.ID)
		if err != nil {
			return err
		}

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

		return nil
	})
}

func (h *dockerKernelHandler) LogsStdoutContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("logsStdoutContainer"),
		oapi.Summary("Get container stdout logs"),
		oapi.Description("Get container stdout logs as a stream."),
	}
}

type LogsStderrContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) LogsStderrContainer() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *LogsStderrContainerParams) error {
		stderr, err := h.dockerService.LogsStderrContainer(params.ID)
		if err != nil {
			return err
		}

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

		return nil
	})
}

func (h *dockerKernelHandler) LogsStderrContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("logsStderrContainer"),
		oapi.Summary("Get container stderr logs"),
		oapi.Description("Get container stderr logs as a stream."),
	}
}

type WaitContainerParams struct {
	ID   string `path:"id"`
	Cond string `path:"cond"`
}

func (h *dockerKernelHandler) WaitContainer() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *WaitContainerParams) error {
		return h.dockerService.WaitContainer(params.ID, types.WaitContainerCondition(params.Cond))
	})
}

func (h *dockerKernelHandler) WaitContainerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("waitContainer"),
		oapi.Summary("Wait container"),
	}
}

type DeleteMountsParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) DeleteMounts() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *DeleteMountsParams) error {
		return h.dockerService.DeleteMounts(params.ID)
	})
}

func (h *dockerKernelHandler) DeleteMountsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("deleteMounts"),
		oapi.Summary("Delete mounts"),
	}
}

type InfoImageParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) InfoImage() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *InfoImageParams) (*types.InfoImageResponse, error) {
		info, err := h.dockerService.InfoImage(params.ID)
		return &info, err
	})
}

func (h *dockerKernelHandler) InfoImageInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("infoImage"),
		oapi.Summary("Get image info"),
	}
}

func (h *dockerKernelHandler) PullImage() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *types.PullImageOptions) error {
		r, err := h.dockerService.PullImage(*params)
		if err != nil {
			return err
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

		return nil
	})
}

func (h *dockerKernelHandler) PullImageInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("pullImage"),
		oapi.Summary("Pull image"),
	}
}

func (h *dockerKernelHandler) BuildImage() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *types.BuildImageOptions) error {
		res, err := h.dockerService.BuildImage(*params)
		if err != nil {
			return err
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

		return nil
	})
}

func (h *dockerKernelHandler) BuildImageInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("buildImage"),
		oapi.Summary("Build image"),
	}
}
