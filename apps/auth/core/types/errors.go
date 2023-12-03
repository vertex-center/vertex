package types

import "github.com/vertex-center/vertex/pkg/router"

const (
	ErrCodeInvalidCredentials         router.ErrCode = "invalid_credentials"
	ErrCodeInvalidToken               router.ErrCode = "invalid_token"
	ErrCodeLoginEmpty                 router.ErrCode = "login_empty"
	ErrCodePasswordEmpty              router.ErrCode = "password_empty"
	ErrCodePasswordLength             router.ErrCode = "password_length"
	ErrCodeFailedToLogout             router.ErrCode = "failed_to_logout"
	ErrCodeFailedToGetUser            router.ErrCode = "failed_to_get_user"
	ErrCodeFailedToPatchUser          router.ErrCode = "failed_to_patch_user"
	ErrCodeFailedToGetUserCredentials router.ErrCode = "failed_to_get_user_credentials"
	ErrCodeFailedToGetUserEmails      router.ErrCode = "failed_to_get_user_emails"
	ErrCodeEmailEmpty                 router.ErrCode = "email_empty"
	ErrCodeInvalidEmail               router.ErrCode = "invalid_email"
	ErrCodeEmailAlreadyExists         router.ErrCode = "email_already_exists"
	ErrCodeFailedToCreateEmail        router.ErrCode = "failed_to_create_email"
	ErrCodeFailedToDeleteEmail        router.ErrCode = "failed_to_delete_email"
)
