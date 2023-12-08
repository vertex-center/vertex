package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
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
		return h.sshService.GetAll()
	})
}

func (h *sshHandler) GetInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getSSHKeys"),
		fizz.Summary("Get all SSH keys"),
	}
}

type AddSSHKeyParams struct {
	AuthorizedKey string `json:"authorized_key"`
	Username      string `json:"username"`
}

func (h *sshHandler) Add() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *AddSSHKeyParams) error {
		return h.sshService.Add(params.AuthorizedKey, params.Username)
	})
}

func (h *sshHandler) AddInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("addSSHKey"),
		fizz.Summary("Add an SSH key"),
	}
}

type DeleteSSHKeyParams struct {
	Fingerprint string `json:"fingerprint"`
	Username    string `json:"username"`
}

func (h *sshHandler) Delete() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *DeleteSSHKeyParams) error {
		return h.sshService.Delete(params.Fingerprint, params.Username)
	})
}

func (h *sshHandler) DeleteInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("deleteSSHKey"),
		fizz.Summary("Delete SSH key"),
	}
}

func (h *sshHandler) GetUsers() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]string, error) {
		return h.sshService.GetUsers()
	})
}

func (h *sshHandler) GetUsersInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getSSHUsers"),
		fizz.Summary("Get all users that can have SSH keys"),
	}
}
