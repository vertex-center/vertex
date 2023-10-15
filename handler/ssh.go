package handler

import (
	"errors"
	"fmt"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/service"
	"github.com/vertex-center/vertex/core/types/api"
	"net/http"

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

type AddSSHKeyBody struct {
	AuthorizedKey string `json:"authorized_key"`
}

func (h *SshHandler) Add(c *router.Context) {
	var body AddSSHKeyBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	err = h.sshService.Add(body.AuthorizedKey)
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

func (h *SshHandler) Delete(c *router.Context) {
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
