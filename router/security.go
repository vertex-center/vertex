package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/types/api"
)

func addSecurityRoutes(r *gin.RouterGroup) {
	r.GET("/ssh", handleGetSSHKey)
	r.POST("/ssh", handleAddSSHKey)
	r.DELETE("/ssh/:fingerprint", handleDeleteSSHKey)
}

// handleGetSSHKey handles the retrieval of the SSH key.
// Errors can be:
//   - failed_to_get_ssh_keys: failed to get the SSH keys.
func handleGetSSHKey(c *gin.Context) {
	keys, err := sshService.GetAll()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToGetSSHKeys,
			Message: fmt.Sprintf("failed to get SSH keys: %v", err),
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
func handleAddSSHKey(c *gin.Context) {
	var body AddSSHKeyBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrFailedToParseBody,
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	err = sshService.Add(body.AuthorizedKey)
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
func handleDeleteSSHKey(c *gin.Context) {
	fingerprint := c.Param("fingerprint")
	if fingerprint == "" {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrInvalidFingerprint,
			Message: "invalid fingerprint",
		})
		return
	}

	err := sshService.Delete(fingerprint)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToDeleteSSHKey,
			Message: fmt.Sprintf("failed to delete SSH key: %v", err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
