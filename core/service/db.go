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

type DbService struct {
	uuid              uuid.UUID
	ctx               *vtypes.VertexContext
	dataConfigAdapter port.DbConfigAdapter
}

func NewDbService(ctx *vtypes.VertexContext, dataConfigAdapter port.DbConfigAdapter) port.DbService {
	s := &DbService{
		uuid:              uuid.New(),
		ctx:               ctx,
		dataConfigAdapter: dataConfigAdapter,
	}
	s.ctx.AddListener(s)
	return s
}

func (s *DbService) GetCurrentDbms() vtypes.DbmsName {
	return s.dataConfigAdapter.GetDBMSName()
}

func (s *DbService) MigrateTo(dbms vtypes.DbmsName) error {
	log.Info("Migrating data to " + string(dbms))

	currentDbms := s.dataConfigAdapter.GetDBMSName()
	if currentDbms == dbms {
		return ErrDbmsAlreadySet
	}

	var err error
	switch dbms {
	case vtypes.DbmsNameSqlite:
		//err = errors.New("sqlite migration not implemented yet")
	case vtypes.DbmsNamePostgres:
		err = s.migrateToPostgres()
	default:
		err = errors.New("unknown dbms: " + string(dbms))
	}

	if err != nil {
		return err
	}

	switch currentDbms {
	case vtypes.DbmsNameSqlite:
		// Nothing to do yet
	case vtypes.DbmsNamePostgres:
		err = s.deletePostgresDB()
	default:
		err = errors.New("unknown dbms: " + string(currentDbms))
	}

	if err != nil {
		return err
	}

	return s.dataConfigAdapter.SetDBMSName(dbms)
}

func (s *DbService) OnEvent(e event.Event) {
	switch e.(type) {
	case vtypes.EventServerStart:
		go func() {
			s.setup()
			s.ctx.DispatchEvent(vtypes.EventServerSetupCompleted{})
		}()
	}
}

func (s *DbService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *DbService) setup() {
	log.Info("Starting Data setup")

	dbms := s.dataConfigAdapter.GetDBMSName()

	var err error
	switch dbms {
	case vtypes.DbmsNameSqlite:
		// Nothing to do yet
	case vtypes.DbmsNamePostgres:
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
