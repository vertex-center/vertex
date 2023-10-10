package instancesapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/instances"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/api"
)

func GetInstance(ctx context.Context, uuid uuid.UUID) (*types.Instance, *api.Error) {
	var inst types.Instance
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s", uuid).
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return &inst, api.HandleError(err, apiError)
}

func DeleteInstance(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s", uuid).
		Delete().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func PatchInstance(ctx context.Context, uuid uuid.UUID, settings types.InstanceSettings) *api.Error {
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s", uuid).
		Patch().
		BodyJSON(&settings).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func StartInstance(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s/start", uuid).
		Patch().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func StopInstance(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s/stop", uuid).
		Patch().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func PatchInstanceEnvironment(ctx context.Context, uuid uuid.UUID, env map[string]string) *api.Error {
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s/environment", uuid).
		Patch().
		BodyJSON(&env).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func GetDocker(ctx context.Context, uuid uuid.UUID) (map[string]any, *api.Error) {
	var info map[string]any
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s/docker", uuid).
		ToJSON(&info).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return info, api.HandleError(err, apiError)
}

func RecreateDocker(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s/docker/recreate", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func GetInstanceLogs(ctx context.Context, uuid uuid.UUID) (string, *api.Error) {
	var logs string
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s/logs", uuid).
		ToJSON(&logs).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return logs, api.HandleError(err, apiError)
}

func UpdateServiceInstance(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s/update/service", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func GetVersions(ctx context.Context, uuid uuid.UUID) ([]string, *api.Error) {
	var versions []string
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Pathf("./instance/%s/versions", uuid).
		ToJSON(&versions).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return versions, api.HandleError(err, apiError)
}

// Helpers

func GetInstanceUUIDParam(c *router.Context) (uuid.UUID, *api.Error) {
	p := c.Param("instance_uuid")
	if p == "" {
		return uuid.UUID{}, &api.Error{
			Code:    api.ErrInstanceUuidMissing,
			Message: "The request was missing the instance UUID.",
		}
	}

	uid, err := uuid.Parse(p)
	if err != nil {
		return uuid.UUID{}, &api.Error{
			Code:    api.ErrInstanceUuidInvalid,
			Message: "The instance UUID is invalid.",
		}
	}

	return uid, nil
}
