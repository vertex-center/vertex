package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"gorm.io/gorm"
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

func (s *MigrationService) migrate(db *gorm.DB) error {
	t := s.getTypes()

	for _, tp := range t {
		err := db.AutoMigrate(tp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MigrationService) getTypes() []interface{} {
	return []interface{}{
		types.User{},
		types.CredentialsArgon2id{},
		types.Token{},
	}
}

func (s *MigrationService) OnEvent(e event.Event) {
	switch e := e.(type) {
	case vtypes.EventDbMigrate:
		err := s.migrate(e.Db)
		if err != nil {
			log.Error(err)
		}
	case vtypes.EventDbCopy:
		e.AddTable(s.getTypes()...)
	}
}

func (s *MigrationService) GetUUID() uuid.UUID {
	return s.uuid
}
