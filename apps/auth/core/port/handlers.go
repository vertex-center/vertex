package port

import (
	"github.com/gin-gonic/gin"
	"github.com/wI2L/fizz"
)

type (
	AuthHandler interface {
		Login() gin.HandlerFunc
		LoginInfo() []fizz.OperationOption

		Register() gin.HandlerFunc
		RegisterInfo() []fizz.OperationOption

		Logout() gin.HandlerFunc
		LogoutInfo() []fizz.OperationOption

		Verify() gin.HandlerFunc
		VerifyInfo() []fizz.OperationOption
	}

	EmailHandler interface {
		CreateCurrentUserEmail() gin.HandlerFunc
		CreateCurrentUserEmailInfo() []fizz.OperationOption

		GetCurrentUserEmails() gin.HandlerFunc
		GetCurrentUserEmailsInfo() []fizz.OperationOption

		DeleteCurrentUserEmail() gin.HandlerFunc
		DeleteCurrentUserEmailInfo() []fizz.OperationOption
	}

	UserHandler interface {
		GetCurrentUser() gin.HandlerFunc
		GetCurrentUserInfo() []fizz.OperationOption

		PatchCurrentUser() gin.HandlerFunc
		PatchCurrentUserInfo() []fizz.OperationOption

		GetCurrentUserCredentials() gin.HandlerFunc
		GetCurrentUserCredentialsInfo() []fizz.OperationOption
	}
)
