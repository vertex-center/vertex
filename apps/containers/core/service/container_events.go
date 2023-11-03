package service

import (
	"github.com/google/uuid"
	vtypes "github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/migration"
	"github.com/vertex-center/vertex/pkg/event/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func (s *ContainerService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *ContainerService) OnEvent(e types.Event) {
	switch e := e.(type) {
	case vtypes.EventServerStart:
		go func() {
			log.Info("post-migration commands", vlog.Any("commands", e.PostMigrationCommands))
			s.LoadAll()
			s.deleteContainersIfNeeded(e.PostMigrationCommands)
			s.StartAll()
			s.ctx.DispatchEvent(vtypes.EventAppReady{
				AppID: "vx-containers",
			})
		}()
	case vtypes.EventServerStop:
		s.StopAll()
	case vtypes.EventServerHardReset:
		s.StopAll()
		s.DeleteAll()
	}
}

func (s *ContainerService) deleteContainersIfNeeded(postMigrationCommands []interface{}) {
	if len(postMigrationCommands) == 0 {
		log.Debug("no post-migration commands", vlog.String("app", "vx-containers"))
		return
	}
	for _, cmd := range postMigrationCommands {
		switch cmd.(type) {
		case migration.CommandRecreateContainers:
			log.Info("post-migration", vlog.String("action", "recreating all containers"))
			containers := s.GetAll()
			for _, c := range containers {
				err := s.containerRunnerService.Delete(c)
				if err != nil {
					log.Error(err)
				}
			}
			return
		}
	}
}
