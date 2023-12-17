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
	return router.Handler(func(ctx *gin.Context) ([]types.PublicKey, error) {
		return h.sshService.GetAll(ctx)
	})
}

type AddSSHKeyParams struct {
	AuthorizedKey string `json:"authorized_key"`
	Username      string `json:"username"`
}

func (h *sshHandler) Add() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *AddSSHKeyParams) error {
		return h.sshService.Add(ctx, params.AuthorizedKey, params.Username)
	})
}

type DeleteSSHKeyParams struct {
	Fingerprint string `json:"fingerprint"`
	Username    string `json:"username"`
}

func (h *sshHandler) Delete() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *DeleteSSHKeyParams) error {
		return h.sshService.Delete(ctx, params.Fingerprint, params.Username)
	})
}

func (h *sshHandler) GetUsers() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) ([]string, error) {
		return h.sshService.GetUsers(ctx)
	})
}
