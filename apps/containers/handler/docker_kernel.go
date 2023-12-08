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
	"github.com/wI2L/fizz"
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

func (h *dockerKernelHandler) GetContainersInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getContainers"),
		fizz.Summary("Get containers"),
	}
}

func (h *dockerKernelHandler) CreateContainer() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *types.CreateContainerOptions) (*types.CreateContainerResponse, error) {
		res, err := h.dockerService.CreateContainer(*params)
		return &res, err
	})
}

func (h *dockerKernelHandler) CreateContainerInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("createContainer"),
		fizz.Summary("Create container"),
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

func (h *dockerKernelHandler) DeleteContainerInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("deleteContainer"),
		fizz.Summary("Delete container"),
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

func (h *dockerKernelHandler) StartContainerInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("startContainer"),
		fizz.Summary("Start container"),
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

func (h *dockerKernelHandler) StopContainerInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("stopContainer"),
		fizz.Summary("Stop container"),
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

func (h *dockerKernelHandler) InfoContainerInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("infoContainer"),
		fizz.Summary("Get container info"),
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

func (h *dockerKernelHandler) LogsStdoutContainerInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("logsStdoutContainer"),
		fizz.Summary("Get container stdout logs"),
		fizz.Description("Get container stdout logs as a stream."),
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

func (h *dockerKernelHandler) LogsStderrContainerInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("logsStderrContainer"),
		fizz.Summary("Get container stderr logs"),
		fizz.Description("Get container stderr logs as a stream."),
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

func (h *dockerKernelHandler) WaitContainerInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("waitContainer"),
		fizz.Summary("Wait container"),
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

func (h *dockerKernelHandler) DeleteMountsInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("deleteMounts"),
		fizz.Summary("Delete mounts"),
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

func (h *dockerKernelHandler) InfoImageInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("infoImage"),
		fizz.Summary("Get image info"),
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

func (h *dockerKernelHandler) PullImageInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("pullImage"),
		fizz.Summary("Pull image"),
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

func (h *dockerKernelHandler) BuildImageInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("buildImage"),
		fizz.Summary("Build image"),
	}
}
