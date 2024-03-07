package types

import (
	"github.com/juju/errors"
	"github.com/vertex-center/uuid"
)

type EnvVariableType string

const (
	EnvVariableTypeString EnvVariableType = "string"
	EnvVariableTypePort   EnvVariableType = "port"
)

var ErrInvalidEnvVariableName = errors.NotValidf("environment variable name")

type (
	EnvVariable struct {
		ID          uuid.UUID       `json:"id"                    db:"id"            example:"7e63ced7-4f4e-4b79-95ca-62930866f7bc"`
		ContainerID uuid.UUID       `json:"container_id"          db:"container_id"  example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		Type        EnvVariableType `json:"type"                  db:"type"          enum:"port"`
		Name        string          `json:"name"                  db:"name"          example:"DB_PORT"`
		DisplayName string          `json:"display_name"          db:"display_name"  example:"Database Port"`
		Value       string          `json:"value"                 db:"value"         example:"5400"`
		Default     *string         `json:"default,omitempty"     db:"default_value" example:"5432"`
		Description *string         `json:"description,omitempty" db:"description"   example:"The server database port"`
		Secret      bool            `json:"secret"                db:"secret"        example:"true"`
	}

	EnvVariableFilters struct {
		ContainerID *uuid.UUID `json:"container_id,omitempty" db:"container_id" example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
	}
)

func (v *EnvVariable) Validate() error {
	if v.Name == "" {
		return ErrInvalidEnvVariableName
	}
	if v.Type == "" {
		v.Type = EnvVariableTypeString
	}
	return nil
}
