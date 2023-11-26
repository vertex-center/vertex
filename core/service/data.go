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
	ErrDbmsAlreadySet           = errors.New("dbms already set")
	ErrPostgresDatabaseNotFound = errors.New("vertex postgres database not found")
)

type DataService struct {
	uuid              uuid.UUID
	ctx               *vtypes.VertexContext
	dataConfigAdapter port.DataConfigAdapter
}

func NewDataService(ctx *vtypes.VertexContext, dataConfigAdapter port.DataConfigAdapter) *DataService {
	s := &DataService{
		uuid:              uuid.New(),
		ctx:               ctx,
		dataConfigAdapter: dataConfigAdapter,
	}
	s.ctx.AddListener(s)
	return s
}

func (s *DataService) GetCurrentDbms() vtypes.DbmsName {
	return s.dataConfigAdapter.GetDBMSName()
}

func (s *DataService) MigrateTo(dbms vtypes.DbmsName) error {
	log.Info("Migrating data to " + string(dbms))

	currentDbms := s.dataConfigAdapter.GetDBMSName()
	if currentDbms == dbms {
		return ErrDbmsAlreadySet
	}

	var err error
	switch dbms {
	case vtypes.DbNameSqlite:
		//err = errors.New("sqlite migration not implemented yet")
	case vtypes.DbNamePostgres:
		err = s.migrateToPostgres()
	default:
		err = errors.New("unknown dbms: " + string(dbms))
	}

	if err != nil {
		return err
	}

	switch currentDbms {
	case vtypes.DbNameSqlite:
		// Nothing to do yet
	case vtypes.DbNamePostgres:
		err = s.deletePostgresDB()
	default:
		err = errors.New("unknown dbms: " + string(currentDbms))
	}

	if err != nil {
		return err
	}

	return s.dataConfigAdapter.SetDBMSName(dbms)
}

func (s *DataService) OnEvent(e event.Event) {
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

func (s *DataService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *DataService) setup() {
	log.Info("Starting Data setup")

	dbms := s.dataConfigAdapter.GetDBMSName()

	var err error
	switch dbms {
	case vtypes.DbNameSqlite:
		// Nothing to do yet
	case vtypes.DbNamePostgres:
		err = s.setupPostgres()
	default:
		log.Error(errors.New("unknown dbms"), vlog.String("dbms", string(dbms)))
	}

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("Data setup completed")
}
