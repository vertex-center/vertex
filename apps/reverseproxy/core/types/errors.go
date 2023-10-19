package types

import (
	"errors"

	"github.com/vertex-center/vertex/pkg/router"
)

const (
	ErrCodeRedirectUuidMissing    router.ErrCode = "redirect_uuid_missing"
	ErrCodeRedirectUuidInvalid    router.ErrCode = "redirect_uuid_invalid"
	ErrCodeFailedToAddRedirect    router.ErrCode = "failed_to_add_redirect"
	ErrCodeFailedToRemoveRedirect router.ErrCode = "failed_to_remove_redirect"
	ErrCodeSourceInvalid          router.ErrCode = "invalid_source_input_to_add_redirect"
	ErrCodeTargetInvalid          router.ErrCode = "invalid_target_input_to_add_redirect"
)

var (
	ErrSourceInvalid = errors.New("source is empty")
	ErrTargetInvalid = errors.New("target is empty")
)
