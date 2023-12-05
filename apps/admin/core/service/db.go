package service

import (
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	vtypes "github.com/vertex-center/vertex/core/types"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/core/types/storage"
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
	uuid uuid.UUID
	ctx  *apptypes.Context
	db   storage.DB
}

func NewDbService(ctx *apptypes.Context, db storage.DB) port.DbService {
	s := &DbService{
		uuid: uuid.New(),
		ctx:  ctx,
		db:   db,
	}
	s.ctx.AddListener(s)
	return s
}

func (s *DbService) GetCurrentDbms() string {
	return s.db.GetDBMSName()
}

func (s *DbService) MigrateTo(dbms string) error {
	log.Info("migrating database", vlog.String("name", dbms))

	currentDbms := s.db.GetDBMSName()
	if currentDbms == dbms {
		return ErrDbmsAlreadySet
	}

	log.Info("setup new dbms", vlog.String("name", dbms))

	var err error
	switch dbms {
	case "sqlite":
		//err = errors.New("sqlite migration not implemented yet")
	case "postgres":
		err = s.setupPostgres()
	default:
		err = errors.New("unknown dbms: " + dbms)
	}
	if err != nil {
		return err
	}

	log.Info("setup new dbms completed", vlog.String("name", dbms))
	log.Info("retrieving connections to previous and next databases", vlog.String("name", dbms))

	prevDb := s.db.DB
	err = s.db.SetDBMSName(dbms)
	if err != nil {
		return err
	}
	err = s.db.Connect()
	if err != nil {
		return err
	}
	nextDb := s.db.DB

	err = s.runMigrations(nextDb)
	if err != nil {
		return err
	}

	log.Info("copying data between databases", vlog.String("from", dbms), vlog.String("to", currentDbms))

	err = s.copyDb(prevDb, nextDb)
	if err != nil {
		return err
	}

	log.Info("copying data between databases completed", vlog.String("from", dbms), vlog.String("to", currentDbms))
	log.Info("deleting previous database", vlog.String("name", currentDbms))

	switch currentDbms {
	case "sqlite":
		err = s.deleteSqliteDB()
	case "postgres":
		err = s.deletePostgresDB()
	default:
		err = errors.New("unknown dbms: " + currentDbms)
	}

	log.Info("deleting previous database completed", vlog.String("name", currentDbms))

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
	dbms := s.db.GetDBMSName()

	log.Info("database setup started", vlog.String("dbms", dbms))

	var err error
	switch dbms {
	case "sqlite":
		// Nothing to do yet
	case "postgres":
		err = s.setupPostgres()
	default:
		log.Error(errors.New("unknown dbms"), vlog.String("dbms", dbms))
	}

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = s.db.Connect()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = s.runMigrations(s.db.DB)
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
