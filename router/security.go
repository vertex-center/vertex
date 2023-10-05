package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/types/api"
)

func addSecurityRoutes(r *router.Group) {
	r.GET("/ssh", handleGetSSHKey)
	r.POST("/ssh", handleAddSSHKey)
	r.DELETE("/ssh/:fingerprint", handleDeleteSSHKey)
}

// handleGetSSHKey handles the retrieval of the SSH key.
// Errors can be:
//   - failed_to_get_ssh_keys: failed to get the SSH keys.
func handleGetSSHKey(c *router.Context) {
	keys, err := sshService.GetAll()
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetSSHKeys,
			PublicMessage:  "Failed to get SSH keys.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, keys)
}

type AddSSHKeyBody struct {
	AuthorizedKey string `json:"authorized_key"`
}

// handleAddSSHKey handles the addition of an SSH key.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - failed_to_add_ssh_key: failed to add the SSH key.
func handleAddSSHKey(c *router.Context) {
	var body AddSSHKeyBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	err = sshService.Add(body.AuthorizedKey)
	if err != nil && errors.Is(err, services.ErrInvalidPublicKey) {
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

// handleDeleteSSHKey handles the deletion of an SSH key.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - failed_to_delete_ssh_key: failed to delete the SSH key.
func handleDeleteSSHKey(c *router.Context) {
	fingerprint := c.Param("fingerprint")
	if fingerprint == "" {
		c.BadRequest(router.Error{
			Code:           api.ErrInvalidFingerprint,
			PublicMessage:  "The request is missing the fingerprint.",
			PrivateMessage: "Field 'fingerprint' is required.",
		})
		return
	}

	err := sshService.Delete(fingerprint)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToDeleteSSHKey,
			PublicMessage:  fmt.Sprintf("Failed to delete SSH key with fingerprint '%s'.", fingerprint),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
