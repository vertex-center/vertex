package handler

import (
	"bufio"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	apierrors "github.com/juju/errors"
	"github.com/vertex-center/vertex/server/apps/containers/core/port"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vertex/server/pkg/router"
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
	return router.Handler(func(ctx *gin.Context) ([]types.DockerContainer, error) {
		return h.dockerService.ListContainers()
	})
}

func (h *dockerKernelHandler) CreateContainer() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *types.CreateDockerContainerOptions) (*types.CreateContainerResponse, error) {
		res, err := h.dockerService.CreateContainer(*params)
		return &res, err
	})
}

type DeleteDockerContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) DeleteContainer() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *DeleteDockerContainerParams) error {
		return h.dockerService.DeleteContainer(params.ID)
	})
}

type StartDockerContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) StartContainer() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *StartDockerContainerParams) error {
		err := h.dockerService.StartContainer(params.ID)
		if err != nil && client.IsErrNotFound(err) {
			return apierrors.NewNotFound(err, "container not found")
		}
		return err
	})
}

type StopDockerContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) StopContainer() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *StopDockerContainerParams) error {
		err := h.dockerService.StopContainer(params.ID)
		if err != nil && client.IsErrNotFound(err) {
			return apierrors.NewNotFound(err, "container not found")
		}
		return err
	})
}

type InfoContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) InfoContainer() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *InfoContainerParams) (*types.InfoContainerResponse, error) {
		info, err := h.dockerService.InfoContainer(params.ID)
		if err != nil && client.IsErrNotFound(err) {
			return nil, apierrors.NewNotFound(err, "container not found")
		}
		return &info, err
	})
}

type LogsStdoutContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) LogsStdoutContainer() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *LogsStdoutContainerParams) error {
		stdout, err := h.dockerService.LogsStdoutContainer(params.ID)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(stdout)

		ctx.Stream(func(w io.Writer) bool {
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

type LogsStderrContainerParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) LogsStderrContainer() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *LogsStderrContainerParams) error {
		stderr, err := h.dockerService.LogsStderrContainer(params.ID)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(stderr)

		ctx.Stream(func(w io.Writer) bool {
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

type WaitContainerParams struct {
	ID   string `path:"id"`
	Cond string `path:"cond"`
}

func (h *dockerKernelHandler) WaitContainer() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *WaitContainerParams) error {
		return h.dockerService.WaitContainer(params.ID, types.WaitContainerCondition(params.Cond))
	})
}

type DeleteMountsParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) DeleteMounts() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *DeleteMountsParams) error {
		return h.dockerService.DeleteMounts(params.ID)
	})
}

type InfoImageParams struct {
	ID string `path:"id"`
}

func (h *dockerKernelHandler) InfoImage() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *InfoImageParams) (*types.InfoImageResponse, error) {
		info, err := h.dockerService.InfoImage(params.ID)
		return &info, err
	})
}

func (h *dockerKernelHandler) PullImage() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *types.PullImageOptions) error {
		r, err := h.dockerService.PullImage(*params)
		if err != nil {
			return err
		}
		defer r.Close()

		scanner := bufio.NewScanner(r)

		ctx.Stream(func(w io.Writer) bool {
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

func (h *dockerKernelHandler) BuildImage() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *types.BuildImageOptions) error {
		res, err := h.dockerService.BuildImage(*params)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		scanner := bufio.NewScanner(res.Body)

		ctx.Stream(func(w io.Writer) bool {
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

type CreateVolumeParams struct {
	Name string `json:"name"`
}

func (h *dockerKernelHandler) CreateVolume() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *CreateVolumeParams) (*volume.Volume, error) {
		res, err := h.dockerService.CreateVolume(params.Name)
		return &res, err
	})
}

type DeleteVolumeParams struct {
	Name string `json:"name"`
}

func (h *dockerKernelHandler) DeleteVolume() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *DeleteVolumeParams) error {
		return h.dockerService.DeleteVolume(params.Name)
	})
}
