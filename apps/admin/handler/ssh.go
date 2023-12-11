package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type sshHandler struct {
	sshService port.SshService
}

func NewSshHandler(sshService port.SshService) port.SshHandler {
	return &sshHandler{
		sshService: sshService,
	}
}

func (h *sshHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]types.PublicKey, error) {
		return h.sshService.GetAll(c)
	})
}

type AddSSHKeyParams struct {
	AuthorizedKey string `json:"authorized_key"`
	Username      string `json:"username"`
}

func (h *sshHandler) Add() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *AddSSHKeyParams) error {
		return h.sshService.Add(c, params.AuthorizedKey, params.Username)
	})
}

type DeleteSSHKeyParams struct {
	Fingerprint string `json:"fingerprint"`
	Username    string `json:"username"`
}

func (h *sshHandler) Delete() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *DeleteSSHKeyParams) error {
		return h.sshService.Delete(c, params.Fingerprint, params.Username)
	})
}

func (h *sshHandler) GetUsers() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]string, error) {
		return h.sshService.GetUsers(c)
	})
}
