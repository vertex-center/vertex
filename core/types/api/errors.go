package api

import "github.com/vertex-center/vertex/pkg/router"

const (
	ErrInternalError router.ErrCode = "internal_error"

	ErrFailedToParseBody router.ErrCode = "failed_to_parse_body"

	ErrFailedToInstallUpdates     router.ErrCode = "failed_to_install_updates"
	ErrAlreadyUpdating            router.ErrCode = "already_updating"
	ErrFailedToFetchLatestVersion router.ErrCode = "failed_to_fetch_latest_version"
	ErrFailedToGetUpdates         router.ErrCode = "failed_to_get_updates"

	ErrFailedToGetSSHKeys   router.ErrCode = "failed_to_get_ssh_keys"
	ErrFailedToAddSSHKey    router.ErrCode = "failed_to_add_ssh_key"
	ErrFailedToDeleteSSHKey router.ErrCode = "failed_to_delete_ssh_key"
	ErrInvalidPublicKey     router.ErrCode = "invalid_public_key"
	ErrInvalidFingerprint   router.ErrCode = "invalid_fingerprint"

	ErrFailedToPatchSettings router.ErrCode = "failed_to_patch_settings"
)
