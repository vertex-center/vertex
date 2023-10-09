package services

import (
	"errors"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/adapter"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type SqlService struct {
	uuid uuid.UUID

	dbms      map[uuid.UUID]types.SqlDBMSAdapterPort
	dbmsMutex *sync.RWMutex

	events types.EventAdapterPort
}

func NewSqlService(eventsAdapter types.EventAdapterPort) SqlService {
	s := SqlService{
		uuid: uuid.New(),

		dbms:      map[uuid.UUID]types.SqlDBMSAdapterPort{},
		dbmsMutex: &sync.RWMutex{},

		events: eventsAdapter,
	}

	s.events.AddListener(&s)

	return s
}

func (s *SqlService) getDbFeature(inst *types.Instance) (types.DatabaseFeature, error) {
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

func (s *SqlService) Get(inst *types.Instance) (types.SqlDBMS, error) {
	db := types.SqlDBMS{}

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

func (s *SqlService) EnvCredentials(inst *types.Instance, user string, pass string) (types.InstanceEnvVariables, error) {
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

func (s *SqlService) onInstanceStart(inst *types.Instance) {
	_, err := s.getDbFeature(inst)
	if err != nil {
		// Not a SQL database
		return
	}

	s.dbmsMutex.Lock()
	defer s.dbmsMutex.Unlock()

	if _, ok := s.dbms[inst.UUID]; ok {
		return
	}

	dbms, err := s.createDbmsAdapter(inst)
	if err != nil {
		return
	}

	s.dbms[inst.UUID] = dbms
}

func (s *SqlService) createDbmsAdapter(inst *types.Instance) (types.SqlDBMSAdapterPort, error) {
	feature, err := s.getDbFeature(inst)
	if err != nil {
		return nil, err
	}

	switch feature.Type {
	case "postgres":
		log.Info("found postgres DBMS", vlog.String("uuid", inst.UUID.String()))
		params := &adapter.SqlDBMSPostgresAdapterParams{
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

		return adapter.NewSqlDBMSPostgresAdapter(params), nil
	default:
		log.Warn("unknown DBMS, generic DBMS used", vlog.String("uuid", inst.UUID.String()))
		return adapter.NewSqlDBMSAdapter(), nil
	}
}

func (s *SqlService) onInstanceStop(uuid uuid.UUID) {
	s.dbmsMutex.Lock()
	defer s.dbmsMutex.Unlock()

	if _, ok := s.dbms[uuid]; !ok {
		return
	}

	delete(s.dbms, uuid)
}

func (s *SqlService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventInstanceStatusChange:
		if e.Status == types.InstanceStatusRunning {
			s.onInstanceStart(&e.Instance)
		} else if e.Status == types.InstanceStatusOff {
			s.onInstanceStop(e.InstanceUUID)
		}
	}
}

func (s *SqlService) GetUUID() uuid.UUID {
	return s.uuid
}
