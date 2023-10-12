package services

import (
	"context"
	"errors"
	"os"
	"strings"

	instancesapi "github.com/vertex-center/vertex/apps/instances/api"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	sqlapi "github.com/vertex-center/vertex/apps/sql/api"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

var (
	ErrPostgresDatabaseNotFound = errors.New("vertex postgres database not found")
)

type SetupService struct {
	ctx *types.VertexContext
}

func NewSetupService(ctx *types.VertexContext) *SetupService {
	return &SetupService{
		ctx: ctx,
	}
}

func (s *SetupService) Setup() error {
	address := config.Current.VertexURL()
	address = strings.TrimPrefix(address, "http://")
	address = strings.TrimPrefix(address, "https://")

	err := net.Wait(address)
	if err != nil {
		return err
	}

	inst, err := s.setupDatabase()
	if err != nil {
		return err
	}

	err = s.startDatabase(inst)
	if err != nil {
		return err
	}

	return nil
}

func (s *SetupService) setupDatabase() (*instancestypes.Instance, error) {
	inst, err := s.getVertexDB()
	if err != nil && !errors.Is(err, ErrPostgresDatabaseNotFound) {
		return nil, err
	}

	if inst != nil {
		log.Info("found vertex postgres instance", vlog.String("uuid", inst.UUID.String()))
		return inst, nil
	}

	err = s.installVertexDB()
	if err != nil {
		return nil, err
	}

	inst, err = s.getVertexDB()
	if err != nil {
		return nil, err
	}

	log.Info("vertex postgres database installed successfully",
		vlog.String("uuid", inst.UUID.String()))

	return inst, nil
}

func (s *SetupService) getVertexDB() (*instancestypes.Instance, error) {
	insts, apiError := instancesapi.GetInstances(context.Background())
	if apiError != nil {
		log.Error(apiError.RouterError())
		os.Exit(1)
	}

	for _, inst := range insts {
		isPostgres, isVertex := false, false
		for _, tag := range inst.Tags {
			if tag == "vertex-postgres-sql" {
				isPostgres = true
			}
			if tag == "vertex" {
				isVertex = true
			}
		}
		if isPostgres && isVertex {
			return inst, nil
		}
	}

	return nil, ErrPostgresDatabaseNotFound
}

func (s *SetupService) installVertexDB() error {
	log.Info("installing vertex postgres database")

	inst, apiError := sqlapi.InstallDBMS(context.Background(), "postgres")
	if apiError != nil {
		return apiError.RouterError()
	}

	inst.Tags = append(inst.Tags, "vertex")
	inst.DisplayName = "Vertex Database"

	apiError = instancesapi.PatchInstance(context.Background(), inst.UUID, inst.InstanceSettings)
	if apiError != nil {
		return apiError.RouterError()
	}

	return nil
}

func (s *SetupService) startDatabase(inst *instancestypes.Instance) error {
	eventsChan := make(chan interface{})
	defer close(eventsChan)

	l := types.NewTempListener(func(e interface{}) {
		switch event := e.(type) {
		case instancestypes.EventInstanceStatusChange:
			if event.InstanceUUID != inst.UUID {
				return
			}
			eventsChan <- event
		}
	})

	s.ctx.AddListener(l)
	defer s.ctx.RemoveListener(l)

	apiError := instancesapi.StartInstance(context.Background(), inst.UUID)
	if apiError != nil {
		return apiError.RouterError()
	}

	errFailedToStart := errors.New("failed to start vertex postgres database")

	for event := range eventsChan {
		switch e := event.(type) {
		case instancestypes.EventInstanceStatusChange:
			if e.Status == instancestypes.InstanceStatusRunning {
				return nil
			} else if e.Status == instancestypes.InstanceStatusOff || e.Status == instancestypes.InstanceStatusError {
				return errFailedToStart
			}
		}
	}

	return errFailedToStart
}
