package containersapi

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/apps/containers/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/api"
)

func GetContainer(ctx context.Context, uuid uuid.UUID) (*types.Container, *api.Error) {
	var inst types.Container
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s", uuid).
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return &inst, api.HandleError(err, apiError)
}

func DeleteContainer(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s", uuid).
		Delete().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func PatchContainer(ctx context.Context, uuid uuid.UUID, settings types.ContainerSettings) *api.Error {
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s", uuid).
		Patch().
		BodyJSON(&settings).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func StartContainer(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s/start", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func StopContainer(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s/stop", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func PatchContainerEnvironment(ctx context.Context, uuid uuid.UUID, env map[string]string) *api.Error {
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s/environment", uuid).
		Patch().
		BodyJSON(&env).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func GetDocker(ctx context.Context, uuid uuid.UUID) (map[string]any, *api.Error) {
	var info map[string]any
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s/docker", uuid).
		ToJSON(&info).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return info, api.HandleError(err, apiError)
}

func RecreateDocker(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s/docker/recreate", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func GetContainerLogs(ctx context.Context, uuid uuid.UUID) (string, *api.Error) {
	var logs string
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s/logs", uuid).
		ToJSON(&logs).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return logs, api.HandleError(err, apiError)
}

func UpdateServiceContainer(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s/update/service", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func GetVersions(ctx context.Context, uuid uuid.UUID) ([]string, *api.Error) {
	var versions []string
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s/versions", uuid).
		ToJSON(&versions).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return versions, api.HandleError(err, apiError)
}

func WaitCondition(ctx context.Context, uuid uuid.UUID, condition container.WaitCondition) *api.Error {
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./container/%s/wait/%s", uuid, condition).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

// Helpers

func GetContainerUUIDParam(c *router.Context) (uuid.UUID, *api.Error) {
	p := c.Param("container_uuid")
	if p == "" {
		return uuid.UUID{}, &api.Error{
			Code:    types.ErrCodeContainerUuidMissing,
			Message: "The request was missing the container UUID.",
		}
	}

	uid, err := uuid.Parse(p)
	if err != nil {
		return uuid.UUID{}, &api.Error{
			Code:    types.ErrCodeContainerUuidInvalid,
			Message: "The container UUID is invalid.",
		}
	}

	return uid, nil
}
