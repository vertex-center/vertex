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
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToGetSSHKeys,
			Message: fmt.Sprintf("failed to get SSH keys: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, keys)
}

// handleAddSSHKey handles the addition of an SSH key.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - failed_to_add_ssh_key: failed to add the SSH key.
func handleAddSSHKeyKernel(c *router.Context) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrFailedToParseBody,
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}
	key := buf.String()

	err = sshKernelService.Add(key)
	if err != nil && errors.Is(err, services.ErrInvalidPublicKey) {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrInvalidPublicKey,
			Message: fmt.Sprintf("error while parsing the public key: %v", err),
		})
		return
	} else if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToAddSSHKey,
			Message: fmt.Sprintf("failed to add SSH key: %v", err),
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
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrInvalidFingerprint,
			Message: "invalid fingerprint",
		})
		return
	}

	err := sshKernelService.Delete(fingerprint)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToDeleteSSHKey,
			Message: fmt.Sprintf("failed to delete SSH key: %v", err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
