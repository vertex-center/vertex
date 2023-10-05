package api

type ErrCode string

const (
	ErrFailedToParseBody ErrCode = "failed_to_parse_body"

	ErrFailedToInstallUpdates ErrCode = "failed_to_install_updates"
	ErrFailedToReloadServices ErrCode = "failed_to_reload_services"

	ErrFailedToListContainers    ErrCode = "failed_to_list_containers"
	ErrFailedToDeleteContainer   ErrCode = "failed_to_delete_container"
	ErrFailedToCreateContainer   ErrCode = "failed_to_create_container"
	ErrFailedToStartContainer    ErrCode = "failed_to_start_container"
	ErrFailedToStopContainer     ErrCode = "failed_to_stop_container"
	ErrFailedToRecreateContainer ErrCode = "failed_to_recreate_container"
	ErrFailedToGetContainerLogs  ErrCode = "failed_to_get_container_logs"
	ErrFailedToWaitContainer     ErrCode = "failed_to_wait_container"
	ErrFailedToGetContainerInfo  ErrCode = "failed_to_get_container_info"
	ErrFailedToGetImageInfo      ErrCode = "failed_to_get_image_info"
	ErrFailedToPullImage         ErrCode = "failed_to_pull_image"
	ErrFailedToBuildImage        ErrCode = "failed_to_build_image"
	ErrContainerNotFound         ErrCode = "container_not_found"

	ErrInstanceUuidInvalid           ErrCode = "instance_uuid_invalid"
	ErrInstanceUuidMissing           ErrCode = "instance_uuid_missing"
	ErrInstanceNotFound              ErrCode = "instance_not_found"
	ErrInstanceAlreadyRunning        ErrCode = "instance_already_running"
	ErrInstanceStillRunning          ErrCode = "instance_still_running"
	ErrInstanceNotRunning            ErrCode = "instance_not_running"
	ErrFailedToGetInstance           ErrCode = "failed_to_get_instance"
	ErrFailedToStartInstance         ErrCode = "failed_to_start_instance"
	ErrFailedToStopInstance          ErrCode = "failed_to_stop_instance"
	ErrFailedToDeleteInstance        ErrCode = "failed_to_delete_instance"
	ErrFailedToGetInstanceLogs       ErrCode = "failed_to_get_logs"
	ErrFailedToUpdateServiceInstance ErrCode = "failed_to_update_service_instance"
	ErrFailedToGetVersions           ErrCode = "failed_to_get_versions"
	ErrFailedToSetLaunchOnStartup    ErrCode = "failed_to_set_launch_on_startup"
	ErrFailedToSetDisplayName        ErrCode = "failed_to_set_display_name"
	ErrFailedToSetDatabase           ErrCode = "failed_to_set_database"
	ErrFailedToSetVersion            ErrCode = "failed_to_set_version"
	ErrFailedToSetEnv                ErrCode = "failed_to_set_env"

	ErrFailedToCheckForUpdates ErrCode = "failed_to_check_for_updates"

	ErrRedirectUuidMissing    ErrCode = "redirect_uuid_missing"
	ErrRedirectUuidInvalid    ErrCode = "redirect_uuid_invalid"
	ErrFailedToAddRedirect    ErrCode = "failed_to_add_redirect"
	ErrFailedToRemoveRedirect ErrCode = "failed_to_remove_redirect"

	ErrFailedToGetSSHKeys   ErrCode = "failed_to_get_ssh_keys"
	ErrFailedToAddSSHKey    ErrCode = "failed_to_add_ssh_key"
	ErrFailedToDeleteSSHKey ErrCode = "failed_to_delete_ssh_key"
	ErrInvalidPublicKey     ErrCode = "invalid_public_key"
	ErrInvalidFingerprint   ErrCode = "invalid_fingerprint"

	ErrServiceNotFound        ErrCode = "service_not_found"
	ErrFailedToInstallService ErrCode = "failed_to_install_service"

	ErrFailedToPatchSettings ErrCode = "failed_to_patch_settings"
)

type Error struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
