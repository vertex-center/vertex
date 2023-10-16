package types

import "github.com/vertex-center/vertex/pkg/router"

const (
	ErrCodeRedirectUuidMissing    router.ErrCode = "redirect_uuid_missing"
	ErrCodeRedirectUuidInvalid    router.ErrCode = "redirect_uuid_invalid"
	ErrCodeFailedToAddRedirect    router.ErrCode = "failed_to_add_redirect"
	ErrCodeFailedToRemoveRedirect router.ErrCode = "failed_to_remove_redirect"
)
