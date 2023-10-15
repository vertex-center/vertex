package api

import "github.com/vertex-center/vertex/pkg/router"

const (
	ErrInternalError router.ErrCode = "internal_error"

	ErrFailedToParseBody router.ErrCode = "failed_to_parse_body"

	ErrFailedToInstallUpdates     router.ErrCode = "failed_to_install_updates"
	ErrAlreadyUpdating            router.ErrCode = "already_updating"
	ErrFailedToFetchLatestVersion router.ErrCode = "failed_to_fetch_latest_version"
	ErrFailedToGetUpdates         router.ErrCode = "failed_to_get_updates"

	ErrFailedToListContainers    router.ErrCode = "failed_to_list_containers"
	ErrFailedToDeleteContainer   router.ErrCode = "failed_to_delete_container"
	ErrFailedToCreateContainer   router.ErrCode = "failed_to_create_container"
	ErrFailedToStartContainer    router.ErrCode = "failed_to_start_container"
	ErrFailedToStopContainer     router.ErrCode = "failed_to_stop_container"
	ErrFailedToRecreateContainer router.ErrCode = "failed_to_recreate_container"
	ErrFailedToGetContainerLogs  router.ErrCode = "failed_to_get_container_logs"
	ErrFailedToWaitContainer     router.ErrCode = "failed_to_wait_container"
	ErrFailedToGetContainerInfo  router.ErrCode = "failed_to_get_container_info"
	ErrFailedToGetImageInfo      router.ErrCode = "failed_to_get_image_info"
	ErrFailedToPullImage         router.ErrCode = "failed_to_pull_image"
	ErrFailedToBuildImage        router.ErrCode = "failed_to_build_image"
	ErrContainerNotFound         router.ErrCode = "container_not_found"

	ErrFailedToGetSSHKeys   router.ErrCode = "failed_to_get_ssh_keys"
	ErrFailedToAddSSHKey    router.ErrCode = "failed_to_add_ssh_key"
	ErrFailedToDeleteSSHKey router.ErrCode = "failed_to_delete_ssh_key"
	ErrInvalidPublicKey     router.ErrCode = "invalid_public_key"
	ErrInvalidFingerprint   router.ErrCode = "invalid_fingerprint"

	ErrFailedToPatchSettings router.ErrCode = "failed_to_patch_settings"
)
