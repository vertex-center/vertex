package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/containers/core/port"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
)

type envHandler struct {
	envService port.EnvService
}

func NewEnvHandler(service port.EnvService) port.EnvHandler {
	return &envHandler{service}
}

type GetContainerEnvParams struct {
	ContainerID uuid.NullUUID `query:"container_id"`
}

func (h *envHandler) GetEnv() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *GetContainerEnvParams) ([]types.EnvVariable, error) {
		filters := types.EnvVariableFilters{}
		if params.ContainerID.Valid {
			filters.ContainerID = &params.ContainerID.UUID
		}
		return h.envService.GetEnvs(ctx, filters)
	}, http.StatusOK)
}

type PatchEnvironmentParams struct {
	EnvID uuid.NullUUID `path:"env_id"`
	types.EnvVariable
}

func (h *envHandler) PatchEnv() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *PatchEnvironmentParams) error {
		params.EnvVariable.ID = params.EnvID.UUID
		return h.envService.PatchEnv(ctx, params.EnvVariable)
	}, http.StatusOK)
}

type DeleteEnvironmentParams struct {
	EnvID uuid.NullUUID `path:"env_id"`
}

func (h *envHandler) DeleteEnv() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *DeleteEnvironmentParams) error {
		return h.envService.DeleteEnv(ctx, params.EnvID.UUID)
	}, http.StatusOK)
}

type CreateEnvironmentParams struct {
	ContainerID uuid.NullUUID `json:"container_id"`
	Type        string        `json:"type"`
	Name        string        `json:"name"`
	Value       string        `json:"value"`
}

func (h *envHandler) CreateEnv() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *CreateEnvironmentParams) error {
		return h.envService.CreateEnv(ctx, types.EnvVariable{
			ContainerID: params.ContainerID.UUID,
			Type:        types.EnvVariableType(params.Type),
			Name:        params.Name,
			Value:       params.Value,
		})
	}, http.StatusCreated)
}
