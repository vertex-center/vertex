package service

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type InstanceServiceService struct {
	adapter types.InstanceServiceAdapterPort
}

func NewInstanceServiceService(adapter types.InstanceServiceAdapterPort) *InstanceServiceService {
	return &InstanceServiceService{
		adapter: adapter,
	}
}

// CheckForUpdate checks if the service of an instance has an update.
// If it has, it sets the instance's ServiceUpdate field.
func (s *InstanceServiceService) CheckForUpdate(inst *types.Instance, latest types.Service) error {
	current := inst.Service
	upToDate := reflect.DeepEqual(latest, current)
	log.Debug("service up-to-date", vlog.Bool("up_to_date", upToDate))
	inst.ServiceUpdate.Available = !upToDate
	return nil
}

// Update updates the service of an instance.
// The service passed is the latest version of the service.
func (s *InstanceServiceService) Update(inst *types.Instance, service types.Service) error {
	if service.Version <= types.MaxSupportedVersion {
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

func (s *InstanceServiceService) Save(inst *types.Instance, service types.Service) error {
	inst.Service = service
	return s.adapter.Save(inst.UUID, service)
}

func (s *InstanceServiceService) Load(uuid uuid.UUID) (types.Service, error) {
	return s.adapter.Load(uuid)
}
