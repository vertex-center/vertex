package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/service"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
)

type SshKernelHandler struct {
	sshService port.SshService
}

func NewSshKernelHandler(sshKernelService port.SshService) port.SshKernelHandler {
	return &SshKernelHandler{
		sshService: sshKernelService,
	}
}

// docapi begin get_ssh_keys_kernel
// docapi method GET
// docapi summary Get all SSH keys
// docapi tags Ssh
// docapi response 200 {[]PublicKey} The list of SSH keys.
// docapi response 500
// docapi end

func (h *SshKernelHandler) Get(c *router.Context) {
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

// docapi begin add_ssh_key_kernel
// docapi method POST
// docapi summary Add an SSH key to the authorized_keys file
// docapi tags Ssh
// docapi body {string} The key to append to the authorized_keys file.
// docapi response 201
// docapi response 400
// docapi response 500
// docapi end

func (h *SshKernelHandler) Add(c *router.Context) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           api.ErrFailedToParseBody,
			PublicMessage:  "Failed to parse request body.",
			PrivateMessage: err.Error(),
		})
		return
	}
	key := buf.String()

	err = h.sshService.Add(key)
	if err != nil && errors.Is(err, service.ErrInvalidPublicKey) {
		c.BadRequest(router.Error{
			Code:           api.ErrInvalidPublicKey,
			PublicMessage:  "Invalid public key.",
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

// docapi begin delete_ssh_key_kernel
// docapi method DELETE
// docapi summary Delete an SSH key from the authorized_keys file
// docapi tags Ssh
// docapi query fingerprint {string} The fingerprint of the SSH key to delete.
// docapi response 204
// docapi response 400
// docapi response 500
// docapi end

func (h *SshKernelHandler) Delete(c *router.Context) {
	fingerprint := c.Param("fingerprint")
	if fingerprint == "" {
		c.BadRequest(router.Error{
			Code:           api.ErrInvalidFingerprint,
			PublicMessage:  "The request is missing the fingerprint.",
			PrivateMessage: "Field 'fingerprint' is required.",
		})
		return
	}

	err := h.sshService.Delete(fingerprint)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToDeleteSSHKey,
			PublicMessage:  fmt.Sprintf("Failed to delete SSH key with fingerprint '%s'.", fingerprint),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// docapi begin get_ssh_users_kernel
// docapi method GET
// docapi summary Get all users that can have SSH keys
// docapi tags Ssh
// docapi response 200 {[]User} The list of users.
// docapi response 500
// docapi end

func (h *SshKernelHandler) GetUsers(c *router.Context) {
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
