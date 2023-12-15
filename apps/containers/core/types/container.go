package types

import (
	"github.com/google/uuid"
	"github.com/juju/errors"
)

const (
	ContainerStatusOff      = "off"
	ContainerStatusBuilding = "building"
	ContainerStatusStarting = "starting"
	ContainerStatusRunning  = "running"
	ContainerStatusStopping = "stopping"
	ContainerStatusError    = "error"
)

var (
	ErrContainerNotFound     = errors.NotFoundf("container")
	ErrContainerStillRunning = errors.New("container still running")
	ErrDatabaseIDNotFound    = errors.NotFoundf("database id")
)

type (
	ContainerID struct{ uuid.UUID }

	Containers []Container
	Container  struct {
		ID              ContainerID `json:"id"                       db:"id"                 example:"1cb8c970-395f-4810-8c9e-e4df35f456e1"`
		ServiceID       string      `json:"service_id"               db:"service_id"         example:"postgres"`
		UserID          uuid.UUID   `json:"user_id"                  db:"user_id"            example:"596ecff2-ca67-4194-947d-59e90920680f"`
		Image           string      `json:"image"                    db:"image"              example:"postgres"`
		ImageTag        string      `json:"image_tag,omitempty"      db:"image_tag"          example:"latest"`
		Status          string      `json:"status"                   db:"status"             example:"running"`
		LaunchOnStartup bool        `json:"launch_on_startup"        db:"launch_on_startup"  example:"true"`
		Name            string      `json:"name"                     db:"name"               example:"Postgres"`
		Description     *string     `json:"description"              db:"description"        example:"An SQL database."`
		Color           *string     `json:"color"                    db:"color"              example:"#336699"`
		Icon            *string     `json:"icon"                     db:"icon"               example:"simpleicons/postgres.svg"`
		Command         *string     `json:"command,omitempty"        db:"command"            example:"tunnel run"`

		Databases     map[string]ContainerID `json:"databases,omitempty"`
		Update        *ContainerUpdate       `json:"update,omitempty"`
		ServiceUpdate ServiceUpdate          `json:"service_update,omitempty"`
	}

	ContainerSearchQuery struct {
		Tags     *[]string `json:"tags,omitempty"`
		Features *[]string `json:"features,omitempty"`
	}

	ContainerUpdate struct {
		CurrentVersion string `json:"current_version"`
		LatestVersion  string `json:"latest_version"`
	}

	DownloadProgress struct {
		ID      string `json:"id"`
		Status  string `json:"status"`
		Current int64  `json:"current,omitempty"`
		Total   int64  `json:"total,omitempty"`
	}
)

func NewContainerID() ContainerID  { return ContainerID{uuid.New()} }
func (ContainerID) Type() string   { return "string" }
func (ContainerID) Format() string { return "uuid" }

func NewContainer(id ContainerID, serviceID string) Container {
	return Container{
		ID:        id,
		ServiceID: serviceID,
		Status:    ContainerStatusOff,
	}
}

func (i *Container) DockerImageVertexName() string { return "vertex_image_" + i.ID.String() }
func (i *Container) DockerContainerName() string   { return "VERTEX_CONTAINER_" + i.ID.String() }

func (i *Container) IsRunning() bool {
	return i.Status != ContainerStatusOff && i.Status != ContainerStatusError
}

func (i *Container) IsBusy() bool {
	return i.Status == ContainerStatusBuilding || i.Status == ContainerStatusStarting || i.Status == ContainerStatusStopping
}

func (i *Container) GetImageNameWithTag() string {
	return i.Image + ":" + i.ImageTag
}
