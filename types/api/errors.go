package api

const (
	ErrFailedToParseBody = "failed_to_parse_body"

	ErrFailedToInstallUpdates = "failed_to_install_updates"
	ErrFailedToReloadServices = "failed_to_reload_services"

	ErrFailedToListContainers    = "failed_to_list_containers"
	ErrFailedToDeleteContainer   = "failed_to_delete_container"
	ErrFailedToCreateContainer   = "failed_to_create_container"
	ErrFailedToStartContainer    = "failed_to_start_container"
	ErrFailedToStopContainer     = "failed_to_stop_container"
	ErrFailedToRecreateContainer = "failed_to_recreate_container"
	ErrFailedToGetContainerLogs  = "failed_to_get_container_logs"
	ErrFailedToWaitContainer     = "failed_to_wait_container"
	ErrFailedToGetContainerInfo  = "failed_to_get_container_info"
	ErrFailedToGetImageInfo      = "failed_to_get_image_info"
	ErrFailedToPullImage         = "failed_to_pull_image"
	ErrFailedToBuildImage        = "failed_to_build_image"
	ErrContainerNotFound         = "container_not_found"

	ErrInstanceUuidInvalid           = "instance_uuid_invalid"
	ErrInstanceUuidMissing           = "instance_uuid_missing"
	ErrInstanceNotFound              = "instance_not_found"
	ErrInstanceAlreadyRunning        = "instance_already_running"
	ErrInstanceStillRunning          = "instance_still_running"
	ErrInstanceNotRunning            = "instance_not_running"
	ErrFailedToGetInstance           = "failed_to_get_instance"
	ErrFailedToStartInstance         = "failed_to_start_instance"
	ErrFailedToStopInstance          = "failed_to_stop_instance"
	ErrFailedToDeleteInstance        = "failed_to_delete_instance"
	ErrFailedToGetInstanceLogs       = "failed_to_get_logs"
	ErrFailedToUpdateServiceInstance = "failed_to_update_service_instance"
	ErrFailedToGetVersions           = "failed_to_get_versions"
	ErrFailedToSetLaunchOnStartup    = "failed_to_set_launch_on_startup"
	ErrFailedToSetDisplayName        = "failed_to_set_display_name"
	ErrFailedToSetDatabase           = "failed_to_set_database"
	ErrFailedToSetVersion            = "failed_to_set_version"
	ErrFailedToSetEnv                = "failed_to_set_env"

	ErrFailedToCheckForUpdates = "failed_to_check_for_updates"

	ErrRedirectUuidMissing    = "redirect_uuid_missing"
	ErrRedirectUuidInvalid    = "redirect_uuid_invalid"
	ErrFailedToAddRedirect    = "failed_to_add_redirect"
	ErrFailedToRemoveRedirect = "failed_to_remove_redirect"

	ErrFailedToGetSSHKeys   = "failed_to_get_ssh_keys"
	ErrFailedToAddSSHKey    = "failed_to_add_ssh_key"
	ErrFailedToDeleteSSHKey = "failed_to_delete_ssh_key"
	ErrInvalidPublicKey     = "invalid_public_key"
	ErrInvalidFingerprint   = "invalid_fingerprint"

	ErrServiceNotFound        = "service_not_found"
	ErrFailedToInstallService = "failed_to_install_service"

	ErrFailedToPatchSettings = "failed_to_patch_settings"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
