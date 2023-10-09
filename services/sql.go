package services

import (
	"errors"

	"github.com/vertex-center/vertex/types"
)

type SqlService struct{}

func NewSqlService() SqlService {
	return SqlService{}
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

func (s *SqlService) Get(inst *types.Instance) (types.SqlDatabase, error) {
	db := types.SqlDatabase{}

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
