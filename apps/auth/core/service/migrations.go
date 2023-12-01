package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	vtypes "github.com/vertex-center/vertex/core/types"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/event"
)

type MigrationService struct {
	uuid uuid.UUID
}

func NewMigrationService(ctx *apptypes.Context) port.MigrationService {
	s := &MigrationService{
		uuid: uuid.New(),
	}
	ctx.AddListener(s)
	return s
}

func (s *MigrationService) getTypes() []string {
	return []string{
		"users",
		"credentials_argon2",
		"credentials_argon2_users",
		"sessions",
	}
}

func (s *MigrationService) OnEvent(e event.Event) error {
	switch e := e.(type) {
	case vtypes.EventDbCopy:
		e.AddTable(s.getTypes()...)
	}
	return nil
}

func (s *MigrationService) GetUUID() uuid.UUID {
	return s.uuid
}
