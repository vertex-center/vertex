package router

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

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
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_list_containers",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, containers)
}

func handleDeleteDockerContainer(c *gin.Context) {
	id := c.Param("id")

	err := dockerKernelService.DeleteContainer(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_delete_container",
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
		_ = c.AbortWithError(http.StatusBadRequest, types.APIError{
			Code:    "failed_to_parse_body",
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	res, err := dockerKernelService.CreateContainer(options)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_create_container",
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
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_start_container",
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
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_stop_container",
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
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_get_container_info",
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
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_get_container_logs",
			Message: err.Error(),
		})
		return
	}

	scanner := bufio.NewScanner(stdout)

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

func handleLogsStderrDockerContainer(c *gin.Context) {
	id := c.Param("id")

	stderr, err := dockerKernelService.LogsStderrContainer(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_get_container_logs",
			Message: err.Error(),
		})
		return
	}

	scanner := bufio.NewScanner(stderr)

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

func handleWaitDockerContainer(c *gin.Context) {
	id := c.Param("id")
	cond := c.Param("cond")

	err := dockerKernelService.WaitContainer(id, types.WaitContainerCondition(cond))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_wait_container",
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
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_get_image_info",
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
		_ = c.AbortWithError(http.StatusBadRequest, types.APIError{
			Code:    "failed_to_parse_body",
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	r, err := dockerKernelService.PullImage(options)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_pull_image",
			Message: err.Error(),
		})
		return
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)

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

func handleBuildDockerImage(c *gin.Context) {
	var options types.BuildImageOptions
	err := c.BindJSON(&options)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, types.APIError{
			Code:    "failed_to_parse_body",
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	res, err := dockerKernelService.BuildImage(options)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_build_image",
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
