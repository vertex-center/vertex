package port

import (
	"github.com/gin-gonic/gin"
)

type (
	AuthHandler interface {
		Login() gin.HandlerFunc
		Register() gin.HandlerFunc
		Logout() gin.HandlerFunc
		Verify() gin.HandlerFunc
	}

	EmailHandler interface {
		CreateCurrentUserEmail() gin.HandlerFunc
		GetCurrentUserEmails() gin.HandlerFunc
		DeleteCurrentUserEmail() gin.HandlerFunc
	}

	UserHandler interface {
		GetCurrentUser() gin.HandlerFunc
		PatchCurrentUser() gin.HandlerFunc
		GetCurrentUserCredentials() gin.HandlerFunc
	}
)
