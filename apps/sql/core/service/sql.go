package service

import (
	"context"
	"sync"

	"github.com/juju/errors"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	sqladapter "github.com/vertex-center/vertex/apps/sql/adapter"
	"github.com/vertex-center/vertex/apps/sql/core/port"
	sqltypes "github.com/vertex-center/vertex/apps/sql/core/types"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/uuid"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vlog"
)

type sqlService struct {
	uuid      uuid.UUID
	dbms      map[uuid.UUID]port.DBMSAdapter
	dbmsMutex *sync.RWMutex
}

func New(ctx *app.Context) port.SqlService {
	s := &sqlService{
		uuid:      uuid.New(),
		dbms:      map[uuid.UUID]port.DBMSAdapter{},
		dbmsMutex: &sync.RWMutex{},
	}
	ctx.AddListener(s)
	return s
}

func (s *sqlService) getDbFeature(c *containerstypes.Container) (containerstypes.DatabaseFeature, error) {
	//if c.Service.Features == nil || c.Service.Features.Databases == nil {
	//	return types.DatabaseFeature{}, errors.New("no databases found")
	//}
	//
	//dbFeatures := *c.Service.Features.Databases
	//for _, dbFeature := range dbFeatures {
	//	if dbFeature.Category == "sql" {
	//		return dbFeature, nil
	//	}
	//}

	return containerstypes.DatabaseFeature{}, errors.New("no sql database found")
}

func (s *sqlService) Get(inst *containerstypes.Container) (sqltypes.DBMS, error) {
	db := sqltypes.DBMS{}
	var err error

	//feature, err := s.getDbFeature(inst)
	//if err != nil {
	//	return db, err
	//}
	//if feature.Username != nil {
	//	db.Username = inst.Env.Get(*feature.Username)
	//}
	//if feature.Password != nil {
	//	db.Password = inst.Env.Get(*feature.Password)
	//}

	s.dbmsMutex.RLock()
	defer s.dbmsMutex.RUnlock()

	if dbms, ok := s.dbms[inst.ID]; ok {
		db.Databases, err = dbms.GetDatabases()
		if err != nil {
			return db, err
		}
	}

	return db, nil
}

func (s *sqlService) Install(ctx context.Context, dbms string) (containerstypes.Container, error) {
	var c containerstypes.Container
	client := containersapi.NewContainersClient(ctx)

	serv, err := client.GetService(ctx, dbms)
	if err != nil {
		return c, err
	}

	c, err = client.InstallService(ctx, serv.ID)
	if err != nil {
		return c, err
	}

	tag, err := client.GetTag(ctx, "Vertex SQL")
	if errors.Is(err, errors.NotFound) {
		tag, err = client.CreateTag(ctx, containerstypes.Tag{
			Name: "Vertex SQL",
		})
		if err != nil {
			return c, err
		}
	}

	err = client.AddContainerTag(ctx, c.ID, tag.ID)
	if err != nil {
		return c, err
	}

	feature, err := s.getDbFeature(&c)
	if err != nil {
		return c, err
	}

	var env containerstypes.EnvVariables
	if feature.Username != nil {
		env.Set(*feature.Username, "postgres")
	}
	if feature.Password != nil {
		env.Set(*feature.Password, "postgres")
	}

	err = client.PatchContainerEnvironment(ctx, c.ID, env)
	return c, err
}

func (s *sqlService) createDbmsAdapter(inst *containerstypes.Container) (port.DBMSAdapter, error) {
	feature, err := s.getDbFeature(inst)
	if err != nil {
		return nil, err
	}

	switch feature.Type {
	case "postgres":
		log.Info("found postgres DBMS", vlog.String("uuid", inst.ID.String()))
		params := &sqladapter.SqlDBMSPostgresAdapterParams{
			Host: config.Current.URL("vertex").String(),
		}

		//params.Port, err = strconv.Atoi(inst.Env.Get(feature.Port))
		//if err != nil {
		//	return nil, err
		//}

		//if feature.Username != nil {
		//	params.Username = inst.Env.Get(*feature.Username)
		//}
		//if feature.Password != nil {
		//	params.Password = inst.Env.Get(*feature.Password)
		//}

		return sqladapter.NewSqlDBMSPostgresAdapter(params), nil
	default:
		log.Warn("unknown DBMS, generic DBMS used", vlog.String("uuid", inst.ID.String()))
		return sqladapter.NewSqlDBMSAdapter(), nil
	}
}
