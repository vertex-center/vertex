package service

import (
	"errors"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	sqladapter "github.com/vertex-center/vertex/apps/sql/adapter"
	"github.com/vertex-center/vertex/apps/sql/core/port"
	sqltypes "github.com/vertex-center/vertex/apps/sql/core/types"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type SqlService struct {
	uuid      uuid.UUID
	dbms      map[uuid.UUID]port.DBMSAdapter
	dbmsMutex *sync.RWMutex
}

func New(ctx *app.Context) port.SqlService {
	s := &SqlService{
		uuid:      uuid.New(),
		dbms:      map[uuid.UUID]port.DBMSAdapter{},
		dbmsMutex: &sync.RWMutex{},
	}
	ctx.AddListener(s)
	return s
}

func (s *SqlService) getDbFeature(inst *types.Container) (types.DatabaseFeature, error) {
	if inst.Service.Features == nil || inst.Service.Features.Databases == nil {
		return types.DatabaseFeature{}, errors.New("no databases found")
	}

	dbFeatures := *inst.Service.Features.Databases
	for _, dbFeature := range dbFeatures {
		if dbFeature.Category == "sql" {
			return dbFeature, nil
		}
	}

	return types.DatabaseFeature{}, errors.New("no sql database found")
}

func (s *SqlService) Get(inst *types.Container) (sqltypes.DBMS, error) {
	db := sqltypes.DBMS{}

	feature, err := s.getDbFeature(inst)
	if err != nil {
		return db, err
	}

	if feature.Username != nil {
		db.Username = inst.Env[*feature.Username]
	}
	if feature.Password != nil {
		db.Password = inst.Env[*feature.Password]
	}

	s.dbmsMutex.RLock()
	defer s.dbmsMutex.RUnlock()

	if dbms, ok := s.dbms[inst.UUID]; ok {
		db.Databases, err = dbms.GetDatabases()
		if err != nil {
			return db, err
		}
	}

	return db, nil
}

func (s *SqlService) EnvCredentials(inst *types.Container, user string, pass string) (types.ContainerEnvVariables, error) {
	env := inst.Env

	feature, err := s.getDbFeature(inst)
	if err != nil {
		return env, err
	}

	if feature.Username != nil {
		env[*feature.Username] = user
	}
	if feature.Password != nil {
		env[*feature.Password] = pass
	}

	return env, nil
}

func (s *SqlService) createDbmsAdapter(inst *types.Container) (port.DBMSAdapter, error) {
	feature, err := s.getDbFeature(inst)
	if err != nil {
		return nil, err
	}

	switch feature.Type {
	case "postgres":
		log.Info("found postgres DBMS", vlog.String("uuid", inst.UUID.String()))
		params := &sqladapter.SqlDBMSPostgresAdapterParams{
			Host: config.Current.Host,
		}

		params.Port, err = strconv.Atoi(inst.Env[feature.Port])
		if err != nil {
			return nil, err
		}

		if feature.Username != nil {
			params.Username = inst.Env[*feature.Username]
		}
		if feature.Password != nil {
			params.Password = inst.Env[*feature.Password]
		}

		return sqladapter.NewSqlDBMSPostgresAdapter(params), nil
	default:
		log.Warn("unknown DBMS, generic DBMS used", vlog.String("uuid", inst.UUID.String()))
		return sqladapter.NewSqlDBMSAdapter(), nil
	}
}
