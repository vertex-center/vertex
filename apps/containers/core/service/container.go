package service

import (
	"bufio"
	"context"
	"encoding/json"
	goerrors "errors"
	"path"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/event"
	vstorage "github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vlog"
)

var (
	ErrContainerAlreadyRunning = errors.New("the container is already running")
	ErrContainerNotRunning     = errors.New("the container is not running")
)

type containerService struct {
	uuid       uuid.UUID
	ctx        *app.Context
	caps       port.CapAdapter // capabilities
	containers port.ContainerAdapter
	vars       port.EnvAdapter
	ports      port.PortAdapter
	volumes    port.VolumeAdapter
	tags       port.TagAdapter
	sysctls    port.SysctlAdapter
	runner     port.RunnerAdapter
	services   port.ServiceAdapter
	logs       port.LogsAdapter
}

func NewContainerService(ctx *app.Context,
	caps port.CapAdapter,
	containers port.ContainerAdapter,
	vars port.EnvAdapter,
	ports port.PortAdapter,
	volumes port.VolumeAdapter,
	tags port.TagAdapter,
	sysctls port.SysctlAdapter,
	runner port.RunnerAdapter,
	services port.ServiceAdapter,
	logs port.LogsAdapter,
) port.ContainerService {
	s := &containerService{
		uuid:       uuid.New(),
		ctx:        ctx,
		caps:       caps,
		containers: containers,
		vars:       vars,
		ports:      ports,
		volumes:    volumes,
		tags:       tags,
		sysctls:    sysctls,
		runner:     runner,
		services:   services,
		logs:       logs,
	}
	s.ctx.AddListener(s)
	return s
}

func (s *containerService) Get(ctx context.Context, id types.ContainerID) (*types.Container, error) {
	return s.containers.GetContainer(ctx, id)
}

func (s *containerService) GetContainers(ctx context.Context) (types.Containers, error) {
	return s.containers.GetContainers(ctx)
}

func (s *containerService) GetTags(ctx context.Context) (types.Tags, error) {
	return s.tags.GetUniqueTags(ctx)
}

// Search returns all containers that match the query.
func (s *containerService) Search(ctx context.Context, query types.ContainerSearchQuery) (types.Containers, error) {
	var containers types.Containers

	all, err := s.containers.GetContainers(ctx)
	if err != nil {
		return containers, err
	}

	for _, c := range all {
		//if query.Features != nil {
		//	if !c.HasFeatureIn(*query.Features) {
		//		continue
		//	}
		//}
		if query.Tags != nil {
			if !c.HasTagIn(*query.Tags) {
				continue
			}
		}
		containers = append(containers, c)
	}

	return containers, nil
}

func (s *containerService) Delete(ctx context.Context, id types.ContainerID) error {
	c, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if c.IsRunning() {
		return types.ErrContainerStillRunning
	}

	err = s.runner.DeleteMounts(ctx, c)
	if err != nil && !errors.Is(err, errors.NotFound) {
		return err
	}

	err = s.runner.DeleteContainer(ctx, c)
	if err != nil && !errors.Is(err, errors.NotFound) {
		return err
	}

	deletes := []func(context.Context, types.ContainerID) error{
		s.caps.DeleteCaps,
		s.ports.DeletePorts,
		s.volumes.DeleteVolumes,
		s.sysctls.DeleteSysctls,
		s.vars.DeleteVariables,
		s.tags.DeleteTags,
		s.containers.DeleteContainer,
	}
	for _, f := range deletes {
		err := f(ctx, id)
		if err != nil {
			return err
		}
	}

	err = s.logs.Unregister(id)
	if err != nil {
		return err
	}

	s.ctx.DispatchEvent(types.EventContainerDeleted{
		ContainerID: c.ID,
		ServiceID:   c.ServiceID,
	})
	s.ctx.DispatchEvent(types.EventContainersChange{})

	return nil
}

func (s *containerService) Start(ctx context.Context, id types.ContainerID) error {
	c, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if c.IsBusy() {
		return nil
	}

	s.ctx.DispatchEvent(types.EventContainerLog{
		ContainerID: c.ID,
		Kind:        types.LogKindOut,
		Message:     types.NewLogLineMessageString("Starting container..."),
	})

	log.Info("starting container", vlog.String("uuid", id.String()))

	if c.IsRunning() {
		s.ctx.DispatchEvent(types.EventContainerLog{
			ContainerID: c.ID,
			Kind:        types.LogKindVertexErr,
			Message:     types.NewLogLineMessageString(ErrContainerAlreadyRunning.Error()),
		})
		return ErrContainerAlreadyRunning
	}

	setStatus := func(status string) {
		s.setStatus(c, status)
	}

	stdout, stderr, err := s.runner.Start(ctx, c, setStatus)
	if err != nil {
		s.setStatus(c, types.ContainerStatusError)
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if strings.HasPrefix(scanner.Text(), "DOWNLOAD") {
				msg := strings.TrimPrefix(scanner.Text(), "DOWNLOAD")

				var downloadProgress types.DownloadProgress
				err := json.Unmarshal([]byte(msg), &downloadProgress)
				if err != nil {
					log.Error(err)
					continue
				}

				s.ctx.DispatchEvent(types.EventContainerLog{
					ContainerID: c.ID,
					Kind:        types.LogKindDownload,
					Message:     types.NewLogLineMessageDownload(&downloadProgress),
				})
				continue
			}

			s.ctx.DispatchEvent(types.EventContainerLog{
				ContainerID: c.ID,
				Kind:        types.LogKindOut,
				Message:     types.NewLogLineMessageString(scanner.Text()),
			})
		}
		if scanner.Err() != nil {
			log.Error(scanner.Err())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			s.ctx.DispatchEvent(types.EventContainerLog{
				ContainerID: c.ID,
				Kind:        types.LogKindErr,
				Message:     types.NewLogLineMessageString(scanner.Text()),
			})
		}
		if scanner.Err() != nil {
			log.Error(scanner.Err())
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := s.WaitStatus(ctx, id, types.ContainerStatusRunning)
		if err != nil {
			log.Error(err)
		}
	}()
	wg.Wait()

	return nil
}

func (s *containerService) StartAll(ctx context.Context) error {
	var ids []types.ContainerID

	// TODO: Retrieve only the containers where LaunchOnStartup is true in the DB.

	all, err := s.containers.GetContainers(ctx)
	if err != nil {
		return err
	}

	for _, inst := range all {
		// vertex containers autostart are managed by the startup service.
		if inst.LaunchOnStartup {
			ids = append(ids, inst.ID)
		}
	}

	// Start them
	for _, id := range ids {
		go func(id types.ContainerID) {
			err = s.Start(ctx, id)
			if err != nil {
				log.Warn("failed to auto-start the container",
					vlog.String("uuid", id.String()),
					vlog.String("reason", err.Error()),
				)
			}
		}(id)
	}

	return nil
}

func (s *containerService) Stop(ctx context.Context, id types.ContainerID) error {
	c, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if c.IsBusy() {
		return nil
	}

	if !c.IsRunning() {
		s.ctx.DispatchEvent(types.EventContainerLog{
			ContainerID: c.ID,
			Kind:        types.LogKindVertexErr,
			Message:     types.NewLogLineMessageString(ErrContainerNotRunning.Error()),
		})
		return ErrContainerNotRunning
	}

	// Log stopped
	s.ctx.DispatchEvent(types.EventContainerLog{
		ContainerID: c.ID,
		Kind:        types.LogKindVertexOut,
		Message:     types.NewLogLineMessageString("Stopping container..."),
	})
	log.Info("stopping container", vlog.String("uuid", c.ID.String()))
	s.setStatus(c, types.ContainerStatusStopping)

	err = s.runner.Stop(ctx, c)
	if err == nil {
		s.ctx.DispatchEvent(types.EventContainerLog{
			ContainerID: c.ID,
			Kind:        types.LogKindVertexOut,
			Message:     types.NewLogLineMessageString("Container stopped."),
		})
		log.Info("container stopped", vlog.String("uuid", c.ID.String()))
		s.setStatus(c, types.ContainerStatusOff)
	} else {
		s.setStatus(c, types.ContainerStatusRunning)
	}

	return err
}

func (s *containerService) StopAll(ctx context.Context) error {
	all, err := s.containers.GetContainers(ctx)
	if err != nil {
		return err
	}

	for _, c := range all {
		err := s.Stop(ctx, c.ID)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (s *containerService) RecreateContainer(ctx context.Context, id types.ContainerID) error {
	c, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if c.IsRunning() {
		err := s.Stop(ctx, id)
		if err != nil {
			return err
		}
	}

	// Make sure to only delete the container!
	// The volumes must be kept here.
	err = s.runner.DeleteContainer(ctx, c)
	if err != nil && !errors.Is(err, errors.NotFound) {
		return err
	}

	return s.Start(ctx, id)
}

func (s *containerService) DeleteAll(ctx context.Context) error {
	all, err := s.containers.GetContainers(ctx)
	if err != nil {
		return err
	}

	var errs []error
	for _, c := range all {
		err := s.Delete(ctx, c.ID)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return goerrors.Join(errs...)
}

func (s *containerService) Install(ctx context.Context, serviceID string) (*types.Container, error) {
	id := types.NewContainerID()

	service, err := s.services.Get(serviceID)
	if err != nil {
		return nil, err
	}

	dir := path.Join(storage.FSPath, id.String())
	if service.Methods.Docker.Clone != nil {
		err := vstorage.CloneRepository(dir, service.Methods.Docker.Clone.Repository)
		if err != nil {
			return nil, err
		}
	}

	c := types.Container{
		ID:              id,
		ServiceID:       serviceID,
		Image:           *service.Methods.Docker.Image,
		ImageTag:        "latest",
		Status:          types.ContainerStatusOff,
		LaunchOnStartup: true,
		Name:            service.Name,
		Description:     &service.Description,
		Color:           service.Color,
		Icon:            service.Icon,
		Command:         service.Methods.Docker.Cmd,
	}

	// Set default env
	env := types.EnvVariables{}
	for _, e := range service.Env {
		env = append(env, types.EnvVariable{
			ContainerID: id,
			Type:        types.EnvVariableType(e.Type),
			Name:        e.Name,
			Value:       e.Default,
			Default:     &e.Default,
			Description: &e.Description,
		})
	}
	c.Env = env

	err = s.vars.CreateVariables(ctx, env)
	if err != nil {
		return nil, err
	}

	// Set default capabilities
	if service.Methods.Docker.Capabilities != nil {
		caps := types.Capabilities{}
		for _, cp := range *service.Methods.Docker.Capabilities {
			caps = append(caps, types.Capability{
				ContainerID: id,
				Name:        cp,
			})
		}
		c.Caps = caps

		err = s.caps.CreateCaps(ctx, caps)
		if err != nil {
			return nil, err
		}
	}

	// Set default ports
	if service.Methods.Docker.Ports != nil {
		ports := types.Ports{}
		for in, out := range *service.Methods.Docker.Ports {
			ports = append(ports, types.Port{
				ContainerID: id,
				In:          in,
				Out:         out,
			})
		}
		c.Ports = ports

		err = s.ports.CreatePorts(ctx, ports)
		if err != nil {
			return nil, err
		}
	}

	// Set default volumes
	if service.Methods.Docker.Volumes != nil {
		volumes := types.Volumes{}
		for in, out := range *service.Methods.Docker.Volumes {
			volumes = append(volumes, types.Volume{
				ContainerID: id,
				In:          in,
				Out:         out,
			})
		}
		c.Volumes = volumes

		err = s.volumes.CreateVolumes(ctx, volumes)
		if err != nil {
			return nil, err
		}
	}

	// Set default sysctls
	if service.Methods.Docker.Sysctls != nil {
		sysctls := types.Sysctls{}
		for name, value := range *service.Methods.Docker.Sysctls {
			sysctls = append(sysctls, types.Sysctl{
				ContainerID: id,
				Name:        name,
				Value:       value,
			})
		}
		c.Sysctls = sysctls

		err = s.sysctls.CreateSysctls(ctx, sysctls)
		if err != nil {
			return nil, err
		}
	}

	err = s.containers.CreateContainer(ctx, c)
	if err != nil {
		return nil, err
	}

	err = s.logs.Register(id)
	if err != nil {
		return nil, err
	}

	c.ResetDefaultEnv(service)
	err = s.vars.CreateVariables(ctx, c.Env)
	if err != nil {
		return nil, err
	}

	s.ctx.DispatchEvent(types.EventContainerCreated{})
	s.ctx.DispatchEvent(types.EventContainersChange{})

	return &c, nil
}

func (s *containerService) CheckForUpdates(ctx context.Context) (types.Containers, error) {
	all, err := s.GetContainers(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range all {
		err := s.runner.CheckForUpdates(ctx, &c)
		if err != nil {
			return all, err
		}
	}

	return all, nil
}

func (s *containerService) SetDatabases(ctx context.Context, c *types.Container, databases map[string]types.ContainerID, options map[string]*types.SetDatabasesOptions) error {
	service, err := s.services.Get(c.ServiceID)
	if err != nil {
		return err
	}

	for db := range databases {
		if _, ok := service.Databases[db]; !ok {
			return types.ErrDatabaseIDNotFound
		}
	}

	c.Databases = databases
	// TODO: Save
	return s.remapDatabaseEnv(ctx, c, options)
}

// remapDatabaseEnv remaps the environment variables of a container.
func (s *containerService) remapDatabaseEnv(ctx context.Context, c *types.Container, options map[string]*types.SetDatabasesOptions) error {
	for databaseID, databaseContainerUUID := range c.Databases {
		db, err := s.Get(ctx, databaseContainerUUID)
		if err != nil {
			return err
		}

		dbService, err := s.services.Get(db.ServiceID)
		if err != nil {
			return err
		}

		cService, err := s.services.Get(c.ServiceID)
		if err != nil {
			return err
		}

		host := config.Current.URL("vertex").String()

		dbEnvNames := (*dbService.Features.Databases)[0]
		iEnvNames := cService.Databases[databaseID].Names

		c.Env.Set(iEnvNames.Host, host)
		c.Env.Set(iEnvNames.Port, db.Env.Get(dbEnvNames.Port))
		if dbEnvNames.Username != nil {
			c.Env.Set(iEnvNames.Username, db.Env.Get(*dbEnvNames.Username))
		}
		if dbEnvNames.Password != nil {
			c.Env.Set(iEnvNames.Password, db.Env.Get(*dbEnvNames.Password))
		}

		if options != nil {
			if modifiedFeature, ok := options[databaseID]; ok {
				if modifiedFeature != nil && modifiedFeature.DatabaseName != nil {
					c.Env.Set(iEnvNames.Database, *modifiedFeature.DatabaseName)
					continue
				}
			}
		}

		if dbEnvNames.DefaultDatabase != nil {
			c.Env.Set(iEnvNames.Database, db.Env.Get(*dbEnvNames.DefaultDatabase))
			continue
		}
	}

	return s.vars.CreateVariables(ctx, c.Env)
}

// SaveEnv saves the environment variables of a container
// and applies them by recreating the container.
func (s *containerService) SaveEnv(ctx context.Context, id types.ContainerID, env types.EnvVariables) error {
	err := s.vars.CreateVariables(ctx, env)
	if err != nil {
		return err
	}
	return s.RecreateContainer(ctx, id)
}

func (s *containerService) GetAllVersions(ctx context.Context, id types.ContainerID, useCache bool) ([]string, error) {
	c, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !useCache || len(c.CacheVersions) == 0 {
		versions, err := s.runner.GetAllVersions(ctx, *c)
		if err != nil {
			return nil, err
		}
		c.CacheVersions = versions
	}

	return c.CacheVersions, nil
}

func (s *containerService) GetContainerInfo(ctx context.Context, id types.ContainerID) (map[string]any, error) {
	c, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.runner.Info(ctx, *c)
}

func (s *containerService) WaitStatus(ctx context.Context, id types.ContainerID, status string) error {
	statusChan := make(chan string)
	defer close(statusChan)

	c, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if c.Status == status {
		return nil
	}

	l := event.NewTempListener(func(e event.Event) error {
		switch e := e.(type) {
		case types.EventContainerStatusChange:
			if e.ContainerUUID != c.ID {
				return nil
			}
			statusChan <- e.Status
		}
		return nil
	})

	s.ctx.AddListener(l)
	defer s.ctx.RemoveListener(l)

	for e := range statusChan {
		if e == status {
			return nil
		}
	}

	return errors.Timeoutf("wait status")
}

func (s *containerService) GetLatestLogs(id types.ContainerID) ([]types.LogLine, error) {
	return s.logs.LoadBuffer(id)
}

func (s *containerService) GetServiceByID(ctx context.Context, id string) (*types.Service, error) {
	serv, err := s.services.Get(id)
	if err != nil {
		return nil, err
	}
	return &serv, nil
}

func (s *containerService) GetServices(ctx context.Context) []types.Service {
	return s.services.GetAll()
}

func (s *containerService) setStatus(c *types.Container, status string) {
	if c.Status == status {
		return
	}

	c.Status = status
	s.ctx.DispatchEvent(types.EventContainersChange{})
	s.ctx.DispatchEvent(types.EventContainerStatusChange{
		ContainerUUID: c.ID,
		ServiceID:     c.ServiceID,
		Container:     *c,
		Name:          c.Name,
		Status:        status,
	})
}
