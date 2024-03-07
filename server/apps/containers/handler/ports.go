package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type portsHandler struct {
	portsService port.PortsService
}

func NewPortsHandler(service port.PortsService) port.PortsHandler {
	return &portsHandler{service}
}

type GetContainerPortsParams struct {
	ContainerID uuid.NullUUID `json:"container_id"`
}

func (h *portsHandler) GetPorts() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *GetContainerPortsParams) (types.Ports, error) {
		filters := types.PortFilters{}
		if params.ContainerID.Valid {
			filters.ContainerID = &params.ContainerID.UUID
		}
		return h.portsService.GetPorts(ctx, filters)
	}, http.StatusOK)
}

type PatchPortParams struct {
	PortID uuid.NullUUID `path:"port_id"`
	In     string        `json:"in"`
	Out    string        `json:"out"`
}

func (h *portsHandler) PatchPort() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *PatchPortParams) error {
		return h.portsService.PatchPort(ctx, types.Port{
			ID:  params.PortID.UUID,
			In:  params.In,
			Out: params.Out,
		})
	}, http.StatusOK)
}

type DeletePortParams struct {
	PortID uuid.NullUUID `path:"port_id"`
}

func (h *portsHandler) DeletePort() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *DeletePortParams) error {
		return h.portsService.DeletePort(ctx, params.PortID.UUID)
	}, http.StatusOK)
}

type CreatePortParams struct {
	ContainerID uuid.NullUUID `json:"container_id"`
	In          string        `json:"in"`
	Out         string        `json:"out"`
}

func (h *portsHandler) CreatePort() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *CreatePortParams) error {
		return h.portsService.CreatePort(ctx, types.Port{
			ContainerID: params.ContainerID.UUID,
			In:          params.In,
			Out:         params.Out,
		})
	}, http.StatusCreated)
}
