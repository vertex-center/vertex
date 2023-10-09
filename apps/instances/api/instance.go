package instancesapi

import (
	"context"

	"github.com/carlmjohnson/requests"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func GetInstance(ctx context.Context, uuid uuid.UUID) (*types.Instance, *types.AppApiError) {
	var inst types.Instance
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s", uuid).
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return &inst, types.HandleError(err, apiError)
}

func DeleteInstance(ctx context.Context, uuid uuid.UUID) *types.AppApiError {
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s", uuid).
		Delete().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return types.HandleError(err, apiError)
}

func PatchInstance(ctx context.Context, uuid uuid.UUID, settings types.InstanceSettings) *types.AppApiError {
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s", uuid).
		Patch().
		BodyJSON(&settings).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return types.HandleError(err, apiError)
}

func StartInstance(ctx context.Context, uuid uuid.UUID) *types.AppApiError {
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s/start", uuid).
		Patch().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return types.HandleError(err, apiError)
}

func StopInstance(ctx context.Context, uuid uuid.UUID) *types.AppApiError {
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s/stop", uuid).
		Patch().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return types.HandleError(err, apiError)
}

func PatchInstanceEnvironment(ctx context.Context, uuid uuid.UUID, env map[string]string) *types.AppApiError {
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s/environment", uuid).
		Patch().
		BodyJSON(&env).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return types.HandleError(err, apiError)
}

func GetDocker(ctx context.Context, uuid uuid.UUID) (map[string]any, *types.AppApiError) {
	var info map[string]any
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s/docker", uuid).
		ToJSON(&info).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return info, types.HandleError(err, apiError)
}

func RecreateDocker(ctx context.Context, uuid uuid.UUID) *types.AppApiError {
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s/docker/recreate", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return types.HandleError(err, apiError)
}

func GetInstanceLogs(ctx context.Context, uuid uuid.UUID) (string, *types.AppApiError) {
	var logs string
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s/logs", uuid).
		ToJSON(&logs).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return logs, types.HandleError(err, apiError)
}

func UpdateServiceInstance(ctx context.Context, uuid uuid.UUID) *types.AppApiError {
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s/update/service", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return types.HandleError(err, apiError)
}

func GetVersions(ctx context.Context, uuid uuid.UUID) ([]string, *types.AppApiError) {
	var versions []string
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s/versions", uuid).
		ToJSON(&versions).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return versions, types.HandleError(err, apiError)
}

// Helpers

func GetInstanceUUIDParam(c *router.Context) (uuid.UUID, *types.AppApiError) {
	p := c.Param("instance_uuid")
	if p == "" {
		return uuid.UUID{}, &types.AppApiError{
			Code:    api.ErrInstanceUuidMissing,
			Message: "The request was missing the instance UUID.",
		}
	}

	uid, err := uuid.Parse(p)
	if err != nil {
		return uuid.UUID{}, &types.AppApiError{
			Code:    api.ErrInstanceUuidInvalid,
			Message: "The instance UUID is invalid.",
		}
	}

	return uid, nil
}
