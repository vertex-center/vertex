package types

type (
	Tags []Tag
	Tag  struct {
		ContainerID ContainerID `json:"container_id" db:"container_id" example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		Tag         string      `json:"tag"          db:"tag"          example:"Vertex SQL"`
	}
)
