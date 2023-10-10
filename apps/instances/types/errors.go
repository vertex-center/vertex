package types

import "github.com/vertex-center/vertex/pkg/router"

const (
	ErrCodeInstanceUuidInvalid           router.ErrCode = "instance_uuid_invalid"
	ErrCodeInstanceUuidMissing           router.ErrCode = "instance_uuid_missing"
	ErrCodeInstanceNotFound              router.ErrCode = "instance_not_found"
	ErrCodeInstanceAlreadyRunning        router.ErrCode = "instance_already_running"
	ErrCodeInstanceStillRunning          router.ErrCode = "instance_still_running"
	ErrCodeInstanceNotRunning            router.ErrCode = "instance_not_running"
	ErrCodeFailedToGetInstance           router.ErrCode = "failed_to_get_instance"
	ErrCodeFailedToStartInstance         router.ErrCode = "failed_to_start_instance"
	ErrCodeFailedToStopInstance          router.ErrCode = "failed_to_stop_instance"
	ErrCodeFailedToDeleteInstance        router.ErrCode = "failed_to_delete_instance"
	ErrCodeFailedToGetInstanceLogs       router.ErrCode = "failed_to_get_logs"
	ErrCodeFailedToUpdateServiceInstance router.ErrCode = "failed_to_update_service_instance"
	ErrCodeFailedToGetVersions           router.ErrCode = "failed_to_get_versions"
	ErrCodeFailedToSetLaunchOnStartup    router.ErrCode = "failed_to_set_launch_on_startup"
	ErrCodeFailedToSetDisplayName        router.ErrCode = "failed_to_set_display_name"
	ErrCodeFailedToSetDatabase           router.ErrCode = "failed_to_set_database"
	ErrCodeFailedToSetVersion            router.ErrCode = "failed_to_set_version"
	ErrCodeFailedToSetTags               router.ErrCode = "failed_to_set_tags"
	ErrCodeFailedToSetEnv                router.ErrCode = "failed_to_set_env"
	ErrCodeFailedToCheckForUpdates       router.ErrCode = "failed_to_check_for_updates"

	ErrCodeServiceIdMissing       router.ErrCode = "service_id_missing"
	ErrCodeServiceNotFound        router.ErrCode = "service_not_found"
	ErrCodeFailedToInstallService router.ErrCode = "failed_to_install_service"
)
