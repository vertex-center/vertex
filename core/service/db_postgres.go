package service

import (
	"context"
	"errors"
	"os"

	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	sqlapi "github.com/vertex-center/vertex/apps/sql/api"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func (s *DbService) setupPostgres() error {
	inst, err := s.getPostgresContainer()
	if err != nil && !errors.Is(err, ErrPostgresDatabaseNotFound) {
		return err
	}

	if inst == nil {
		err = s.installPostgresDB()
		if err != nil {
			return err
		}

		inst, err = s.getPostgresContainer()
		if err != nil {
			return err
		}

		log.Info("vertex postgres database installed successfully", vlog.String("uuid", inst.UUID.String()))
	} else {
		log.Info("found vertex postgres container", vlog.String("uuid", inst.UUID.String()))

		client := containersapi.NewContainersClient()
		apiError := client.StartContainer(context.Background(), inst.UUID)
		if apiError != nil {
			return apiError.RouterError()
		}
	}

	log.Info("postgres database is ready")

	return nil
}

func (s *DbService) deletePostgresDB() error {
	inst, err := s.getPostgresContainer()
	if err != nil {
		return err
	}

	client := containersapi.NewContainersClient()

	apiError := client.StopContainer(context.Background(), inst.UUID)
	if apiError != nil {
		return apiError.RouterError()
	}

	apiError = client.DeleteContainer(context.Background(), inst.UUID)
	if apiError != nil {
		return apiError.RouterError()
	}

	return nil
}

func (s *DbService) getPostgresContainer() (*types.Container, error) {
	client := containersapi.NewContainersClient()

	insts, apiError := client.GetContainers(context.Background())
	if apiError != nil {
		log.Error(apiError.RouterError())
		os.Exit(1)
	}

	for _, inst := range insts {
		isDatabase, isVertex, isPostgres := false, false, false
		if inst.Service.Features != nil && inst.Service.Features.Databases != nil {
			for _, db := range *inst.Service.Features.Databases {
				if db.Type == "postgres" {
					isPostgres = true
				}
			}
		}
		for _, tag := range inst.Tags {
			if tag == "Vertex SQL" {
				isDatabase = true
			}
			if tag == "Vertex Internal" {
				isVertex = true
			}
		}
		if isDatabase && isVertex && isPostgres {
			return inst, nil
		}
	}

	return nil, ErrPostgresDatabaseNotFound
}

func (s *DbService) installPostgresDB() error {
	log.Info("installing vertex postgres database")

	sqlClient := sqlapi.NewSqlClient()

	inst, apiError := sqlClient.InstallDBMS(context.Background(), "postgres")
	if apiError != nil {
		return apiError.RouterError()
	}

	inst.Tags = append(inst.Tags, "Vertex Internal")
	inst.DisplayName = "Vertex Database"

	client := containersapi.NewContainersClient()

	apiError = client.PatchContainer(context.Background(), inst.UUID, inst.ContainerSettings)
	if apiError != nil {
		return apiError.RouterError()
	}

	return nil
}

func (s *DbService) migrateToPostgres() error {
	err := s.setupPostgres()
	if err != nil {
		return err
	}

	// TODO: Migrate data

	return nil
}
