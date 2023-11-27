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
	log.Info("migrating database", vlog.String("name", string(dbms)))

	currentDbms := s.dataConfigAdapter.GetDBMSName()
	if currentDbms == dbms {
		return ErrDbmsAlreadySet
	}

	log.Info("setup new dbms", vlog.String("name", string(dbms)))

	var err error
	switch dbms {
	case vtypes.DbmsNameSqlite:
		//err = errors.New("sqlite migration not implemented yet")
	case vtypes.DbmsNamePostgres:
		err = s.setupPostgres()
	default:
		err = errors.New("unknown dbms: " + string(dbms))
	}
	if err != nil {
		return err
	}

	log.Info("setup new dbms completed", vlog.String("name", string(dbms)))
	log.Info("retrieving connections to previous and next databases", vlog.String("name", string(dbms)))

	prevDb := s.dataConfigAdapter.Get()
	err = s.dataConfigAdapter.SetDBMSName(dbms)
	if err != nil {
		return err
	}
	err = s.dataConfigAdapter.Connect()
	if err != nil {
		return err
	}
	nextDb := s.dataConfigAdapter.Get()

	log.Info("copying data between databases", vlog.String("from", string(dbms)), vlog.String("to", string(currentDbms)))

	err = s.copyDb(prevDb, nextDb)
	if err != nil {
		return err
	}

	log.Info("copying data between databases completed", vlog.String("from", string(dbms)), vlog.String("to", string(currentDbms)))
	log.Info("deleting previous database", vlog.String("name", string(currentDbms)))

	switch currentDbms {
	case vtypes.DbmsNameSqlite:
		err = s.deleteSqliteDB()
	case vtypes.DbmsNamePostgres:
		err = s.deletePostgresDB()
	default:
		err = errors.New("unknown dbms: " + string(currentDbms))
	}

	log.Info("deleting previous database completed", vlog.String("name", string(currentDbms)))

	return err
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
	dbms := s.dataConfigAdapter.GetDBMSName()

	log.Info("database setup started", vlog.String("dbms", string(dbms)))

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

	err = s.dataConfigAdapter.Connect()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("database setup completed")
}
