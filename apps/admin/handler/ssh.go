package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/service"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
)

type SshHandler struct {
	sshService port.SshService
}

func NewSshHandler(sshService port.SshService) port.SshHandler {
	return &SshHandler{
		sshService: sshService,
	}
}

// docapi begin get_ssh_keys
// docapi method GET
// docapi summary Get all SSH keys
// docapi tags SSH
// docapi response 200 {[]PublicKey} The list of SSH keys.
// docapi response 500
// docapi end

func (h *SshHandler) Get(c *router.Context) {
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

// docapi begin add_ssh_key
// docapi method POST
// docapi summary Add an SSH key
// docapi tags SSH
// docapi body {AddSSHKeyBody} The SSH key to add.
// docapi response 201
// docapi response 400
// docapi response 500
// docapi end

type AddSSHKeyBody struct {
	AuthorizedKey string `json:"authorized_key"`
	Username      string `json:"username"`
}

func (h *SshHandler) Add(c *router.Context) {
	var body AddSSHKeyBody
	err := c.ParseBody(&body)
	if err != nil {
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
		c.BadRequest(router.Error{
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

	c.Status(http.StatusCreated)
}

// docapi begin delete_ssh_key
// docapi method DELETE
// docapi summary Delete SSH key
// docapi tags SSH
// docapi query fingerprint {DeleteSSHKeyBody} The fingerprint of the SSH key to delete.
// docapi response 204
// docapi response 400
// docapi response 500
// docapi end

type DeleteSSHKeyBody struct {
	Fingerprint string `json:"fingerprint"`
	Username    string `json:"username"`
}

func (h *SshHandler) Delete(c *router.Context) {
	var body DeleteSSHKeyBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	err = h.sshService.Delete(body.Fingerprint, body.Username)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToDeleteSSHKey,
			PublicMessage:  fmt.Sprintf("Failed to delete SSH key with fingerprint '%s' of the user '%s'.", body.Fingerprint, body.Username),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// docapi begin get_ssh_users
// docapi method GET
// docapi summary Get all users that can have SSH keys
// docapi tags SSH
// docapi response 200 {[]User} The list of users.
// docapi response 500
// docapi end

func (h *SshHandler) GetUsers(c *router.Context) {
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
