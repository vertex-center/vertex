package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/service"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
	"github.com/vertex-center/vertex/pkg/user"
)

type sshKernelHandler struct {
	sshService port.SshKernelService
}

func NewSshKernelHandler(sshKernelService port.SshKernelService) port.SshKernelHandler {
	return &sshKernelHandler{
		sshService: sshKernelService,
	}
}

func (h *sshKernelHandler) Get(c *router.Context) {
	keys, err := h.sshService.GetAll()
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetSSHKeys,
			PublicMessage:  "Failed to get SSH keys.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(keys)
}

func (h *sshKernelHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get all SSH keys"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel([]types.PublicKey{}),
		),
	}
}

func (h *sshKernelHandler) Add(c *router.Context) {
	var body AddSSHKeyBody
	err := c.ParseBody(&body)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           api.ErrFailedToParseBody,
			PublicMessage:  "Failed to parse request body.",
			PrivateMessage: err.Error(),
		})
		return
	}

	err = h.sshService.Add(body.AuthorizedKey, body.Username)
	if err != nil && errors.Is(err, service.ErrInvalidPublicKey) {
		c.BadRequest(router.Error{
			Code:           api.ErrInvalidPublicKey,
			PublicMessage:  "Invalid public key.",
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil && errors.Is(err, service.ErrUserNotFound) {
		c.NotFound(router.Error{
			Code:           api.ErrUserNotFound,
			PublicMessage:  "User not found.",
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToAddSSHKey,
			PublicMessage:  "Failed to add SSH key.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.Created()
}

func (h *sshKernelHandler) AddInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Add an SSH key to the authorized_keys file"),
		oapi.Response(http.StatusCreated),
	}
}

func (h *sshKernelHandler) Delete(c *router.Context) {
	var body DeleteSSHKeyBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	err = h.sshService.Delete(body.Fingerprint, body.Username)
	if err != nil && errors.Is(err, service.ErrUserNotFound) {
		c.NotFound(router.Error{
			Code:           api.ErrUserNotFound,
			PublicMessage:  "User not found.",
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToDeleteSSHKey,
			PublicMessage:  fmt.Sprintf("Failed to delete SSH key with fingerprint '%s' of the user '%s'.", body.Fingerprint, body.Username),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *sshKernelHandler) DeleteInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Delete an SSH key from the authorized_keys file"),
		oapi.Response(http.StatusNoContent),
	}
}

func (h *sshKernelHandler) GetUsers(c *router.Context) {
	users, err := h.sshService.GetUsers()
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetSshUsers,
			PublicMessage:  "Failed to get ssh users.",
			PrivateMessage: err.Error(),
		})
		return
	}
	c.JSON(users)
}

func (h *sshKernelHandler) GetUsersInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get all users that can have SSH keys"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel([]user.User{}),
		),
	}
}
