package api

import "github.com/vertex-center/vertex/pkg/router"

const (
	ErrFailedToParseBody router.ErrCode = "failed_to_parse_body"

	ErrFailedToInstallUpdates router.ErrCode = "failed_to_install_updates"
	ErrFailedToReloadServices router.ErrCode = "failed_to_reload_services"

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

	ErrInstanceUuidInvalid           router.ErrCode = "instance_uuid_invalid"
	ErrInstanceUuidMissing           router.ErrCode = "instance_uuid_missing"
	ErrInstanceNotFound              router.ErrCode = "instance_not_found"
	ErrInstanceAlreadyRunning        router.ErrCode = "instance_already_running"
	ErrInstanceStillRunning          router.ErrCode = "instance_still_running"
	ErrInstanceNotRunning            router.ErrCode = "instance_not_running"
	ErrFailedToGetInstance           router.ErrCode = "failed_to_get_instance"
	ErrFailedToStartInstance         router.ErrCode = "failed_to_start_instance"
	ErrFailedToStopInstance          router.ErrCode = "failed_to_stop_instance"
	ErrFailedToDeleteInstance        router.ErrCode = "failed_to_delete_instance"
	ErrFailedToGetInstanceLogs       router.ErrCode = "failed_to_get_logs"
	ErrFailedToUpdateServiceInstance router.ErrCode = "failed_to_update_service_instance"
	ErrFailedToGetVersions           router.ErrCode = "failed_to_get_versions"
	ErrFailedToSetLaunchOnStartup    router.ErrCode = "failed_to_set_launch_on_startup"
	ErrFailedToSetDisplayName        router.ErrCode = "failed_to_set_display_name"
	ErrFailedToSetDatabase           router.ErrCode = "failed_to_set_database"
	ErrFailedToSetVersion            router.ErrCode = "failed_to_set_version"
	ErrFailedToSetTags               router.ErrCode = "failed_to_set_tags"
	ErrFailedToSetEnv                router.ErrCode = "failed_to_set_env"

	ErrFailedToCheckForUpdates router.ErrCode = "failed_to_check_for_updates"

	ErrRedirectUuidMissing    router.ErrCode = "redirect_uuid_missing"
	ErrRedirectUuidInvalid    router.ErrCode = "redirect_uuid_invalid"
	ErrFailedToAddRedirect    router.ErrCode = "failed_to_add_redirect"
	ErrFailedToRemoveRedirect router.ErrCode = "failed_to_remove_redirect"

	ErrFailedToGetSSHKeys   router.ErrCode = "failed_to_get_ssh_keys"
	ErrFailedToAddSSHKey    router.ErrCode = "failed_to_add_ssh_key"
	ErrFailedToDeleteSSHKey router.ErrCode = "failed_to_delete_ssh_key"
	ErrInvalidPublicKey     router.ErrCode = "invalid_public_key"
	ErrInvalidFingerprint   router.ErrCode = "invalid_fingerprint"

	ErrServiceIdMissing       router.ErrCode = "service_id_missing"
	ErrServiceNotFound        router.ErrCode = "service_not_found"
	ErrFailedToInstallService router.ErrCode = "failed_to_install_service"

	ErrFailedToPatchSettings router.ErrCode = "failed_to_patch_settings"

	ErrCollectorNotFound                router.ErrCode = "collector_not_found"
	ErrVisualizerNotFound               router.ErrCode = "visualizer_not_found"
	ErrFailedToConfigureMetricsInstance router.ErrCode = "failed_to_configure_metrics_instance"
	ErrFailedToGetMetrics               router.ErrCode = "failed_to_get_metrics"

	ErrFailedToConfigureTunnelInstance router.ErrCode = "failed_to_configure_tunnel_instance"

	ErrSQLDatabaseNotFound                  router.ErrCode = "sql_database_not_found"
	ErrFailedToConfigureSQLDatabaseInstance router.ErrCode = "failed_to_configure_sql_database_instance"
)
