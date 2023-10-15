package service

import (
	"github.com/vertex-center/vertex/apps/containers/core/port"
	types2 "github.com/vertex-center/vertex/apps/containers/core/types"
	"reflect"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type ContainerServiceService struct {
	adapter port.ContainerServiceAdapter
}

func NewContainerServiceService(adapter port.ContainerServiceAdapter) *ContainerServiceService {
	return &ContainerServiceService{
		adapter: adapter,
	}
}

// CheckForUpdate checks if the service of an container has an update.
// If it has, it sets the container's ServiceUpdate field.
func (s *ContainerServiceService) CheckForUpdate(inst *types2.Container, latest types2.Service) error {
	current := inst.Service
	upToDate := reflect.DeepEqual(latest, current)
	log.Debug("service up-to-date", vlog.Bool("up_to_date", upToDate))
	inst.ServiceUpdate.Available = !upToDate
	return nil
}

// Update updates the service of an container.
// The service passed is the latest version of the service.
func (s *ContainerServiceService) Update(inst *types2.Container, service types2.Service) error {
	if service.Version <= types2.MaxSupportedVersion {
		log.Info("service version is outdated, upgrading.",
			vlog.String("uuid", inst.UUID.String()),
			vlog.Int("old_version", int(inst.Service.Version)),
			vlog.Int("new_version", int(service.Version)),
		)
		err := s.Save(inst, service)
		if err != nil {
			return err
		}
	} else {
		log.Info("service version is not supported, skipping.",
			vlog.String("uuid", inst.UUID.String()),
			vlog.Int("version", int(service.Version)),
		)
	}

	return nil
}

func (s *ContainerServiceService) Save(inst *types2.Container, service types2.Service) error {
	inst.Service = service
	return s.adapter.Save(inst.UUID, service)
}

func (s *ContainerServiceService) Load(uuid uuid.UUID) (types2.Service, error) {
	return s.adapter.Load(uuid)
}
