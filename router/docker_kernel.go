package router

import (
	"bufio"
	"fmt"
	"io"

	"github.com/docker/docker/client"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/api"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
)

func addDockerKernelRoutes(r *router.Group) {
	r.GET("/containers", handleListDockerContainers)
	r.POST("/container", handleCreateDockerContainer)
	r.DELETE("/container/:id", handleDeleteDockerContainer)
	r.POST("/container/:id/start", handleStartDockerContainer)
	r.POST("/container/:id/stop", handleStopDockerContainer)
	r.GET("/container/:id/info", handleInfoDockerContainer)
	r.GET("/container/:id/logs/stdout", handleLogsStdoutDockerContainer)
	r.GET("/container/:id/logs/stderr", handleLogsStderrDockerContainer)
	r.GET("/container/:id/wait/:cond", handleWaitDockerContainer)
	r.GET("/image/:id/info", handleInfoDockerImage)
	r.POST("/image/pull", handlePullDockerImage)
	r.POST("/image/build", handleBuildDockerImage)
}

func handleListDockerContainers(c *router.Context) {
	containers, err := dockerKernelService.ListContainers()
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

func handleDeleteDockerContainer(c *router.Context) {
	id := c.Param("id")

	err := dockerKernelService.DeleteContainer(id)
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

func handleCreateDockerContainer(c *router.Context) {
	var options types.CreateContainerOptions
	err := c.ParseBody(&options)
	if err != nil {
		return
	}

	res, err := dockerKernelService.CreateContainer(options)
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

func handleStartDockerContainer(c *router.Context) {
	id := c.Param("id")

	err := dockerKernelService.StartContainer(id)
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

func handleStopDockerContainer(c *router.Context) {
	id := c.Param("id")

	err := dockerKernelService.StopContainer(id)
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

func handleInfoDockerContainer(c *router.Context) {
	id := c.Param("id")

	info, err := dockerKernelService.InfoContainer(id)
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

func handleLogsStdoutDockerContainer(c *router.Context) {
	id := c.Param("id")

	stdout, err := dockerKernelService.LogsStdoutContainer(id)
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

func handleLogsStderrDockerContainer(c *router.Context) {
	id := c.Param("id")

	stderr, err := dockerKernelService.LogsStderrContainer(id)
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

func handleWaitDockerContainer(c *router.Context) {
	id := c.Param("id")
	cond := c.Param("cond")

	err := dockerKernelService.WaitContainer(id, types.WaitContainerCondition(cond))
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

func handleInfoDockerImage(c *router.Context) {
	id := c.Param("id")

	info, err := dockerKernelService.InfoImage(id)
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

func handlePullDockerImage(c *router.Context) {
	var options types.PullImageOptions
	err := c.ParseBody(&options)
	if err != nil {
		return
	}

	r, err := dockerKernelService.PullImage(options)
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

func handleBuildDockerImage(c *router.Context) {
	var options types.BuildImageOptions
	err := c.ParseBody(&options)
	if err != nil {
		return
	}

	res, err := dockerKernelService.BuildImage(options)
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
