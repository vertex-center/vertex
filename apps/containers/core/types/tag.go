package types

type (
	Tags []Tag
	Tag  struct {
		ContainerUUID ContainerID `json:"container_uuid" db:"container_uuid" example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		Tag           string      `json:"tag"            db:"tag"            example:"Vertex SQL"`
	}
)
