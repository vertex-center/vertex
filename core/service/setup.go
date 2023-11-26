package service

import (
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/core/port"
	vtypes "github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

var (
	ErrPostgresDatabaseNotFound = errors.New("vertex postgres database not found")
)

type SetupService struct {
	uuid              uuid.UUID
	ctx               *vtypes.VertexContext
	dataConfigAdapter port.DataConfigAdapter
}

func NewSetupService(ctx *vtypes.VertexContext, dataConfigAdapter port.DataConfigAdapter) *SetupService {
	s := &SetupService{
		uuid:              uuid.New(),
		ctx:               ctx,
		dataConfigAdapter: dataConfigAdapter,
	}
	s.ctx.AddListener(s)
	return s
}

func (s *SetupService) OnEvent(e event.Event) {
	switch e := e.(type) {
	case vtypes.EventAppReady:
		if e.AppID != "vx-containers" {
			return
		}
		go func() {
			s.setup()
			s.ctx.DispatchEvent(vtypes.EventServerSetupCompleted{})
		}()
	}
}

func (s *SetupService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *SetupService) setup() {
	log.Info("Starting vertex setup")

	dbms := s.dataConfigAdapter.GetDBMSName()

	var err error
	switch dbms {
	case vtypes.DbNamePostgres:
		err = s.setupPostgres()
	case vtypes.DbNameSqlite:
		// Nothing to do yet
	default:
		log.Error(errors.New("unknown dbms"), vlog.String("dbms", string(dbms)))
	}

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("Vertex setup completed")
}
