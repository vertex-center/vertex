package types

import "github.com/vertex-center/vertex/pkg/router"

const (
	ErrCodeInvalidCredentials router.ErrCode = "invalid_credentials"
	ErrCodeInvalidToken       router.ErrCode = "invalid_token"
	ErrCodeLoginEmpty         router.ErrCode = "login_empty"
	ErrCodePasswordEmpty      router.ErrCode = "password_empty"
	ErrCodePasswordLength     router.ErrCode = "password_length"
	ErrCodeFailedToLogout     router.ErrCode = "failed_to_logout"
	ErrCodeFailedToGetUser    router.ErrCode = "failed_to_get_user"
)
