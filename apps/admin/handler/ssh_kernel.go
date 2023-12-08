package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/user"
	"github.com/wI2L/fizz"
)

type sshKernelHandler struct {
	sshService port.SshKernelService
}

func NewSshKernelHandler(sshKernelService port.SshKernelService) port.SshKernelHandler {
	return &sshKernelHandler{
		sshService: sshKernelService,
	}
}

func (h *sshKernelHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]types.PublicKey, error) {
		return h.sshService.GetAll()
	})
}

func (h *sshKernelHandler) GetInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getSSHKeys"),
		fizz.Summary("Get all SSH keys"),
	}
}

func (h *sshKernelHandler) Add() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *AddSSHKeyParams) error {
		return h.sshService.Add(params.AuthorizedKey, params.Username)
	})
}

func (h *sshKernelHandler) AddInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("addSSHKey"),
		fizz.Summary("Add an SSH key to the authorized_keys file"),
	}
}

func (h *sshKernelHandler) Delete() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *DeleteSSHKeyParams) error {
		return h.sshService.Delete(params.Fingerprint, params.Username)
	})
}

func (h *sshKernelHandler) DeleteInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("deleteSSHKey"),
		fizz.Summary("Delete an SSH key from the authorized_keys file"),
	}
}

func (h *sshKernelHandler) GetUsers() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]user.User, error) {
		return h.sshService.GetUsers()
	})
}

func (h *sshKernelHandler) GetUsersInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getSSHUsers"),
		fizz.Summary("Get all users that can have SSH keys"),
	}
}
