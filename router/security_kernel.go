package router

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/types/api"
)

func addSecurityKernelRoutes(r *router.Group) {
	r.GET("/ssh", handleGetSSHKeyKernel)
	r.POST("/ssh", handleAddSSHKeyKernel)
	r.DELETE("/ssh/:fingerprint", handleDeleteSSHKeyKernel)
}

// handleGetSSHKey handles the retrieval of the SSH key.
// Errors can be:
//   - failed_to_get_ssh_keys: failed to get the SSH keys.
func handleGetSSHKeyKernel(c *router.Context) {
	keys, err := sshKernelService.GetAll()
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

// handleAddSSHKey handles the addition of an SSH key.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - failed_to_add_ssh_key: failed to add the SSH key.
func handleAddSSHKeyKernel(c *router.Context) {
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

	err = sshKernelService.Add(key)
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
func handleDeleteSSHKeyKernel(c *router.Context) {
	fingerprint := c.Param("fingerprint")
	if fingerprint == "" {
		c.BadRequest(router.Error{
			Code:           api.ErrInvalidFingerprint,
			PublicMessage:  "The request is missing the fingerprint.",
			PrivateMessage: "Field 'fingerprint' is required.",
		})
		return
	}

	err := sshKernelService.Delete(fingerprint)
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
