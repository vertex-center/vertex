package router

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/vertex-center/vertex/types/api"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
)

func addDockerKernelRoutes(r *gin.RouterGroup) {
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

func handleListDockerContainers(c *gin.Context) {
	containers, err := dockerKernelService.ListContainers()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToListContainers,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, containers)
}

func handleDeleteDockerContainer(c *gin.Context) {
	id := c.Param("id")

	err := dockerKernelService.DeleteContainer(id)
	if err != nil && client.IsErrNotFound(err) {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrContainerNotFound,
			Message: err.Error(),
		})
		return
	} else if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToDeleteContainer,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func handleCreateDockerContainer(c *gin.Context) {
	var options types.CreateContainerOptions
	err := c.BindJSON(&options)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrFailedToParseBody,
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	res, err := dockerKernelService.CreateContainer(options)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToCreateContainer,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func handleStartDockerContainer(c *gin.Context) {
	id := c.Param("id")

	err := dockerKernelService.StartContainer(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToStartContainer,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func handleStopDockerContainer(c *gin.Context) {
	id := c.Param("id")

	err := dockerKernelService.StopContainer(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToStopContainer,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func handleInfoDockerContainer(c *gin.Context) {
	id := c.Param("id")

	info, err := dockerKernelService.InfoContainer(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToGetContainerInfo,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

func handleLogsStdoutDockerContainer(c *gin.Context) {
	id := c.Param("id")

	stdout, err := dockerKernelService.LogsStdoutContainer(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToGetContainerLogs,
			Message: err.Error(),
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

func handleLogsStderrDockerContainer(c *gin.Context) {
	id := c.Param("id")

	stderr, err := dockerKernelService.LogsStderrContainer(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToGetContainerLogs,
			Message: err.Error(),
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

func handleWaitDockerContainer(c *gin.Context) {
	id := c.Param("id")
	cond := c.Param("cond")

	err := dockerKernelService.WaitContainer(id, types.WaitContainerCondition(cond))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToWaitContainer,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func handleInfoDockerImage(c *gin.Context) {
	id := c.Param("id")

	info, err := dockerKernelService.InfoImage(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToGetImageInfo,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

func handlePullDockerImage(c *gin.Context) {
	var options types.PullImageOptions
	err := c.BindJSON(&options)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrFailedToParseBody,
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	r, err := dockerKernelService.PullImage(options)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToPullImage,
			Message: err.Error(),
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

func handleBuildDockerImage(c *gin.Context) {
	var options types.BuildImageOptions
	err := c.BindJSON(&options)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrFailedToParseBody,
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	res, err := dockerKernelService.BuildImage(options)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToBuildImage,
			Message: err.Error(),
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
