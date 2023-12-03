package service

import (
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/database"
	"github.com/vertex-center/vertex/database/migration"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/vsql"
	"github.com/vertex-center/vlog"
)

var (
	ErrDbmsAlreadySet           = errors.New("dbms already set")
	ErrPostgresDatabaseNotFound = errors.New("vertex postgres database not found")
)

type DbService struct {
	uuid              uuid.UUID
	ctx               *apptypes.Context
	dataConfigAdapter port.DbAdapter
}

func NewDbService(ctx *apptypes.Context, dataConfigAdapter port.DbAdapter) port.DbService {
	s := &DbService{
		uuid:              uuid.New(),
		ctx:               ctx,
		dataConfigAdapter: dataConfigAdapter,
	}
	s.ctx.AddListener(s)
	return s
}

func (s *DbService) GetCurrentDbms() types.DbmsName {
	return s.dataConfigAdapter.GetDBMSName()
}

func (s *DbService) MigrateTo(dbms types.DbmsName) error {
	log.Info("migrating database", vlog.String("name", string(dbms)))

	currentDbms := s.dataConfigAdapter.GetDBMSName()
	if currentDbms == dbms {
		return ErrDbmsAlreadySet
	}

	log.Info("setup new dbms", vlog.String("name", string(dbms)))

	var err error
	switch dbms {
	case types.DbmsNameSqlite:
		//err = errors.New("sqlite migration not implemented yet")
	case types.DbmsNamePostgres:
		err = s.setupPostgres()
	default:
		err = errors.New("unknown dbms: " + string(dbms))
	}
	if err != nil {
		return err
	}

	log.Info("setup new dbms completed", vlog.String("name", string(dbms)))
	log.Info("retrieving connections to previous and next databases", vlog.String("name", string(dbms)))

	prevDb := s.dataConfigAdapter.Get().DB
	err = s.dataConfigAdapter.SetDBMSName(dbms)
	if err != nil {
		return err
	}
	err = s.dataConfigAdapter.Connect()
	if err != nil {
		return err
	}
	nextDb := s.dataConfigAdapter.Get().DB

	err = s.runMigrations(nextDb)
	if err != nil {
		return err
	}

	log.Info("copying data between databases", vlog.String("from", string(dbms)), vlog.String("to", string(currentDbms)))

	err = s.copyDb(prevDb, nextDb)
	if err != nil {
		return err
	}

	log.Info("copying data between databases completed", vlog.String("from", string(dbms)), vlog.String("to", string(currentDbms)))
	log.Info("deleting previous database", vlog.String("name", string(currentDbms)))

	switch currentDbms {
	case types.DbmsNameSqlite:
		err = s.deleteSqliteDB()
	case types.DbmsNamePostgres:
		err = s.deletePostgresDB()
	default:
		err = errors.New("unknown dbms: " + string(currentDbms))
	}

	log.Info("deleting previous database completed", vlog.String("name", string(currentDbms)))

	return err
}

func (s *DbService) OnEvent(e event.Event) error {
	switch e.(type) {
	// This needs containers and sql to work
	case vtypes.EventAllAppsReady:
		go func() {
			s.setup()
			s.ctx.DispatchEvent(vtypes.EventServerSetupCompleted{})
		}()
	}
	return nil
}

func (s *DbService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *DbService) setup() {
	dbms := s.dataConfigAdapter.GetDBMSName()

	log.Info("database setup started", vlog.String("dbms", string(dbms)))

	var err error
	switch dbms {
	case types.DbmsNameSqlite:
		// Nothing to do yet
	case types.DbmsNamePostgres:
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

	err = s.runMigrations(s.dataConfigAdapter.Get().DB)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("database setup completed")
}

func (s *DbService) runMigrations(db *sqlx.DB) error {
	var version int
	err := db.Get(&version, "SELECT version FROM migrations LIMIT 1")
	if err != nil {
		return s.createSchemas(db)
	}
	return vsql.Migrate(migration.Migrations, db, version)
}

func (s *DbService) createSchemas(db *sqlx.DB) error {
	vsqlDriver := vsql.DriverFromName(db.DriverName())
	_, err := db.Exec(database.GetSchema(vsqlDriver))
	return err
}
