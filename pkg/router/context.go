package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

func (c *Context) AbortWithError(statusCode int, err Error) {
	c.Header("Content-Type", "application/json")
	_ = c.Context.AbortWithError(statusCode, err)
}

// 200

func (c *Context) JSON(data interface{}) {
	c.Context.JSON(http.StatusOK, data)
}

func (c *Context) Created() {
	c.Context.Status(http.StatusCreated)
}

func (c *Context) OK() {
	c.Context.Status(http.StatusNoContent)
}

// 300

func (c *Context) NotModified() {
	c.Context.Status(http.StatusNotModified)
}

// 400

func (c *Context) BadRequest(err Error) {
	c.AbortWithError(http.StatusBadRequest, err)
}

func (c *Context) Unauthorized(err Error) {
	c.AbortWithError(http.StatusUnauthorized, err)
}

func (c *Context) NotFound(err Error) {
	c.AbortWithError(http.StatusNotFound, err)
}

func (c *Context) Conflict(err Error) {
	c.AbortWithError(http.StatusConflict, err)
}

func (c *Context) Unprocessable(err Error) {
	c.AbortWithError(http.StatusUnprocessableEntity, err)
}

// 500

func (c *Context) Abort(err Error) {
	c.AbortWithError(http.StatusInternalServerError, err)
}

func (c *Context) AbortWithCode(code int, err Error) {
	c.AbortWithError(code, err)
}

func (c *Context) ParseBody(obj interface{}) error {
	err := c.BindJSON(obj)
	if err != nil {
		c.BadRequest(Error{
			Code:           ErrFailedToParseBody,
			PublicMessage:  "Failed to parse the request.",
			PrivateMessage: err.Error(),
		})
		return err
	}
	return nil
}
