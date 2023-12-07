package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	AuthHandler interface {
		Login() gin.HandlerFunc
		LoginInfo() []oapi.Info

		Register() gin.HandlerFunc
		RegisterInfo() []oapi.Info

		Logout() gin.HandlerFunc
		LogoutInfo() []oapi.Info

		Verify() gin.HandlerFunc
		VerifyInfo() []oapi.Info
	}

	EmailHandler interface {
		CreateCurrentUserEmail() gin.HandlerFunc
		CreateCurrentUserEmailInfo() []oapi.Info

		GetCurrentUserEmails() gin.HandlerFunc
		GetCurrentUserEmailsInfo() []oapi.Info

		DeleteCurrentUserEmail() gin.HandlerFunc
		DeleteCurrentUserEmailInfo() []oapi.Info
	}

	UserHandler interface {
		GetCurrentUser() gin.HandlerFunc
		GetCurrentUserInfo() []oapi.Info

		PatchCurrentUser() gin.HandlerFunc
		PatchCurrentUserInfo() []oapi.Info

		GetCurrentUserCredentials() gin.HandlerFunc
		GetCurrentUserCredentialsInfo() []oapi.Info
	}
)
