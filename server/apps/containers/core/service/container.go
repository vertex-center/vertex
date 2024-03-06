package service

import (
	"bufio"
	"context"
	"encoding/json"
	goerrors "errors"
	"path"
	"strings"
	"sync"

	"github.com/juju/errors"
	"github.com/vertex-center/uuid"
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
	templates  port.TemplateAdapter
	logs       port.LogsAdapter

	cacheImageTags map[string][]string
	mu             sync.RWMutex
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
	services port.TemplateAdapter,
	logs port.LogsAdapter,
) port.ContainerService {
	s := &containerService{
		uuid:           uuid.New(),
		ctx:            ctx,
		caps:           caps,
		containers:     containers,
		vars:           vars,
		ports:          ports,
		volumes:        volumes,
		tags:           tags,
		sysctls:        sysctls,
		runner:         runner,
		templates:      services,
		logs:           logs,
		cacheImageTags: make(map[string][]string),
	}
	s.ctx.AddListener(s)
	return s
}

func (s *containerService) Get(ctx context.Context, id uuid.UUID) (*types.Container, error) {
	return s.containers.GetContainer(ctx, id)
}

func (s *containerService) GetContainers(ctx context.Context) (types.Containers, error) {
	return s.containers.GetContainers(ctx)
}

func (s *containerService) GetContainersWithFilters(ctx context.Context, filters types.ContainerFilters) (types.Containers, error) {
	return s.containers.GetContainersWithFilters(ctx, filters)
}

func (s *containerService) CreateContainer(ctx context.Context, opts types.CreateContainerOptions) (*types.Container, error) {
	opts.FillDefaults()

	var (
		id          = uuid.New()
		image       = opts.Image
		imageTag    = opts.ImageTag
		name        = opts.Image
		description *string
		color       *string
		icon        *string
		cmd         *string

		env     []types.TemplateEnv
		caps    []string
		ports   []types.TemplatePort
		volumes = map[string]string{}
		sysctls = map[string]string{}
	)

	if opts.TemplateID != nil {
		template, err := s.templates.Get(*opts.TemplateID)
		if err != nil {
			return nil, err
		}

		image = template.Methods.Docker.Image
		name = &template.Name
		description = &template.Description
		color = template.Color
		icon = template.Icon
		cmd = template.Methods.Docker.Cmd

		env = template.Env
		ports = template.Ports
		if template.Methods.Docker.Capabilities != nil {
			caps = *template.Methods.Docker.Capabilities
		}
		if template.Methods.Docker.Volumes != nil {
			volumes = *template.Methods.Docker.Volumes
		}
		if template.Methods.Docker.Sysctls != nil {
			sysctls = *template.Methods.Docker.Sysctls
		}

		dir := path.Join(storage.FSPath, id.String())
		if template.Methods.Docker.Clone != nil {
			err := vstorage.CloneRepository(dir, template.Methods.Docker.Clone.Repository)
			if err != nil {
				return nil, err
			}
		}
	}

	c := types.Container{
		ID:              id,
		TemplateID:      opts.TemplateID,
		Image:           *image,
		ImageTag:        *imageTag,
		Status:          types.ContainerStatusOff,
		LaunchOnStartup: true,
		Name:            *name,
		Description:     description,
		Color:           color,
		Icon:            icon,
		Command:         cmd,
	}

	err := s.containers.CreateContainer(ctx, c)
	if err != nil {
		return nil, err
	}

	// Set default env
	for _, e := range env {
		err = s.vars.CreateVariable(ctx, types.EnvVariable{
			ID:          uuid.New(),
			ContainerID: id,
			Type:        types.EnvVariableType(e.Type),
			Name:        e.Name,
			DisplayName: e.DisplayName,
			Value:       e.Default,
			Default:     &e.Default,
			Description: &e.Description,
			Secret:      e.Secret != nil && *e.Secret,
		})
		if err != nil {
			return nil, err
		}
	}

	// Set default capabilities
	for _, cp := range caps {
		err = s.caps.CreateCap(ctx, types.Capability{
			ID:          uuid.New(),
			ContainerID: id,
			Name:        cp,
		})
		if err != nil {
			return nil, err
		}
	}

	// Set default ports
	for _, p := range ports {
		err = s.ports.CreatePort(ctx, types.Port{
			ID:          uuid.New(),
			ContainerID: id,
			In:          p.Port,
			Out:         p.Port,
		})
		if err != nil {
			return nil, err
		}
	}

	// Set default volumes
	for out, in := range volumes {
		tp := types.VolumeTypeBind
		if !strings.Contains(out, "/") {
			tp = types.VolumeTypeVolume
			out = "VERTEX_VOLUME_" + id.String() + "_" + out
		}
		err = s.volumes.CreateVolume(ctx, types.Volume{
			ID:          uuid.New(),
			ContainerID: id,
			Type:        tp,
			In:          in,
			Out:         out,
		})
		if err != nil {
			return nil, err
		}
	}

	// Set default sysctls
	for name, value := range sysctls {
		err = s.sysctls.CreateSysctl(ctx, types.Sysctl{
			ID:          uuid.New(),
			ContainerID: id,
			Name:        name,
			Value:       value,
		})
		if err != nil {
			return nil, err
		}
	}

	err = s.logs.Register(id)
	if err != nil {
		return nil, err
	}

	s.ctx.DispatchEvent(types.EventContainerCreated{})
	s.ctx.DispatchEvent(types.EventContainersChange{})

	return &c, nil
}

func (s *containerService) Delete(ctx context.Context, id uuid.UUID) error {
	c, err := s.containers.GetContainer(ctx, id)
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

	volumes, err := s.volumes.GetContainerVolumes(ctx, id)
	if err != nil {
		return err
	}
	var volNames []string
	for _, v := range volumes {
		volNames = append(volNames, v.Out)
	}

	err = s.runner.DeleteContainer(ctx, c, volNames)
	if err != nil && !errors.Is(err, errors.NotFound) {
		return err
	}

	deletes := []func(context.Context, uuid.UUID) error{
		s.caps.DeleteContainerCaps,
		s.ports.DeleteContainerPorts,
		s.volumes.DeleteContainerVolumes,
		s.sysctls.DeleteContainerSysctls,
		s.vars.DeleteContainerVariables,
		s.containers.DeleteTags,
		s.containers.DeleteContainer,
	}
	for _, f := range deletes {
		err := f(ctx, id)
		if err != nil {
			return err
		}
	}

	err = s.logs.Unregister(id)
	if err != nil && !errors.Is(err, errors.NotFound) {
		return err
	}

	s.ctx.DispatchEvent(types.EventContainerDeleted{
		ContainerID: id,
	})
	s.ctx.DispatchEvent(types.EventContainersChange{})

	return nil
}

func (s *containerService) UpdateContainer(ctx context.Context, id uuid.UUID, c types.Container) error {
	c.ID = id
	return s.containers.UpdateContainer(ctx, c)
}

func (s *containerService) Start(ctx context.Context, id uuid.UUID) error {
	c, err := s.containers.GetContainer(ctx, id)
	if err != nil {
		return err
	}

	if c.IsBusy() {
		return nil
	}

	if !s.logs.Exists(id) {
		err = s.logs.Register(id)
		if err != nil {
			return err
		}
	}

	s.ctx.DispatchEvent(types.EventContainerLog{
		ContainerID: id,
		Kind:        types.LogKindOut,
		Message:     types.NewLogLineMessageString("Starting container..."),
	})

	log.Info("starting container", vlog.String("id", id.String()))

	if c.IsRunning() {
		s.ctx.DispatchEvent(types.EventContainerLog{
			ContainerID: id,
			Kind:        types.LogKindVertexErr,
			Message:     types.NewLogLineMessageString(ErrContainerAlreadyRunning.Error()),
		})
		return ErrContainerAlreadyRunning
	}

	setStatus := func(status string) {
		s.setStatus(c, status)
	}

	ports, err := s.ports.GetContainerPorts(ctx, id)
	if err != nil {
		s.setStatus(c, types.ContainerStatusError)
		return err
	}

	volumes, err := s.volumes.GetContainerVolumes(ctx, id)
	if err != nil {
		s.setStatus(c, types.ContainerStatusError)
		return err
	}

	env, err := s.vars.GetContainerVariables(ctx, id)
	if err != nil {
		s.setStatus(c, types.ContainerStatusError)
		return err
	}

	caps, err := s.caps.GetContainerCaps(ctx, id)
	if err != nil {
		s.setStatus(c, types.ContainerStatusError)
		return err
	}

	sysctls, err := s.sysctls.GetContainerSysctls(ctx, id)
	if err != nil {
		s.setStatus(c, types.ContainerStatusError)
		return err
	}

	stdout, stderr, err := s.runner.Start(ctx, c, ports, volumes, env, caps, sysctls, setStatus)
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
					ContainerID: id,
					Kind:        types.LogKindDownload,
					Message:     types.NewLogLineMessageDownload(&downloadProgress),
				})
				continue
			}

			s.ctx.DispatchEvent(types.EventContainerLog{
				ContainerID: id,
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
				ContainerID: id,
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
	var ids []uuid.UUID

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
		go func(id uuid.UUID) {
			err = s.Start(ctx, id)
			if err != nil {
				log.Warn("failed to auto-start the container",
					vlog.String("id", id.String()),
					vlog.String("reason", err.Error()),
				)
			}
		}(id)
	}

	return nil
}

func (s *containerService) Stop(ctx context.Context, id uuid.UUID) error {
	c, err := s.containers.GetContainer(ctx, id)
	if err != nil {
		return err
	}

	if c.IsBusy() {
		return nil
	}

	if !c.IsRunning() {
		s.ctx.DispatchEvent(types.EventContainerLog{
			ContainerID: id,
			Kind:        types.LogKindVertexErr,
			Message:     types.NewLogLineMessageString(ErrContainerNotRunning.Error()),
		})
		return ErrContainerNotRunning
	}

	// Log stopped
	s.ctx.DispatchEvent(types.EventContainerLog{
		ContainerID: id,
		Kind:        types.LogKindVertexOut,
		Message:     types.NewLogLineMessageString("Stopping container..."),
	})
	log.Info("stopping container", vlog.String("id", id.String()))
	s.setStatus(c, types.ContainerStatusStopping)

	err = s.runner.Stop(ctx, c)
	if err == nil {
		s.ctx.DispatchEvent(types.EventContainerLog{
			ContainerID: id,
			Kind:        types.LogKindVertexOut,
			Message:     types.NewLogLineMessageString("Container stopped."),
		})
		log.Info("container stopped", vlog.String("id", id.String()))
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

func (s *containerService) AddContainerTag(ctx context.Context, id uuid.UUID, tagID uuid.UUID) error {
	return s.containers.AddTag(ctx, id, tagID)
}

func (s *containerService) RecreateContainer(ctx context.Context, id uuid.UUID) error {
	c, err := s.containers.GetContainer(ctx, id)
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
	err = s.runner.DeleteContainer(ctx, c, []string{})
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

func (s *containerService) SetDatabases(ctx context.Context, c *types.Container, databases map[string]uuid.UUID, options map[string]*types.SetDatabasesOptions) error {
	if c.TemplateID == nil {
		return nil
	}

	template, err := s.templates.Get(*c.TemplateID)
	if err != nil {
		return err
	}

	for db := range databases {
		if _, ok := template.Databases[db]; !ok {
			return types.ErrDatabaseIDNotFound
		}
	}

	c.Databases = databases
	// TODO: Save
	return s.remapDatabaseEnv(ctx, c, options)
}

func (s *containerService) GetContainerEnv(ctx context.Context, id uuid.UUID) (types.EnvVariables, error) {
	return s.vars.GetContainerVariables(ctx, id)
}

// remapDatabaseEnv remaps the environment variables of a container.
func (s *containerService) remapDatabaseEnv(ctx context.Context, c *types.Container, options map[string]*types.SetDatabasesOptions) error {
	for databaseID, databaseContainerID := range c.Databases {
		db, err := s.containers.GetContainer(ctx, databaseContainerID)
		if err != nil {
			return err
		}

		dbService, err := s.templates.Get(*db.TemplateID)
		if err != nil {
			return err
		}

		cService, err := s.templates.Get(*c.TemplateID)
		if err != nil {
			return err
		}

		host := config.Current.Addr("vertex").String()

		dbEnvNames := (*dbService.Features.Databases)[0]
		cEnvNames := cService.Databases[databaseID].Names

		dbVars, err := s.vars.GetContainerVariables(ctx, db.ID)
		if err != nil {
			return err
		}

		err = s.vars.UpdateContainerVariableByName(ctx, types.EnvVariable{
			ContainerID: c.ID,
			Name:        cEnvNames.Host,
			Value:       host,
		})
		if err != nil {
			return err
		}
		err = s.vars.UpdateContainerVariableByName(ctx, types.EnvVariable{
			ContainerID: c.ID,
			Name:        cEnvNames.Port,
			Value:       dbVars.Get(dbEnvNames.Port),
		})
		if err != nil {
			return err
		}

		if dbEnvNames.Username != nil {
			err = s.vars.UpdateContainerVariableByName(ctx, types.EnvVariable{
				ContainerID: c.ID,
				Name:        cEnvNames.Username,
				Value:       dbVars.Get(*dbEnvNames.Username),
			})
			if err != nil {
				return err
			}
		}
		if dbEnvNames.Password != nil {
			err = s.vars.UpdateContainerVariableByName(ctx, types.EnvVariable{
				ContainerID: c.ID,
				Name:        cEnvNames.Password,
				Value:       dbVars.Get(*dbEnvNames.Password),
			})
			if err != nil {
				return err
			}
		}

		if options != nil {
			if modifiedFeature, ok := options[databaseID]; ok {
				if modifiedFeature != nil && modifiedFeature.DatabaseName != nil {
					err = s.vars.UpdateContainerVariableByName(ctx, types.EnvVariable{
						ContainerID: c.ID,
						Name:        cEnvNames.Database,
						Value:       *modifiedFeature.DatabaseName,
					})
					if err != nil {
						return err
					}
					continue
				}
			}
		}

		if dbEnvNames.DefaultDatabase != nil {
			err = s.vars.UpdateContainerVariableByName(ctx, types.EnvVariable{
				ContainerID: c.ID,
				Name:        cEnvNames.Database,
				Value:       dbVars.Get(*dbEnvNames.DefaultDatabase),
			})
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

// SaveEnv saves the environment variables of a container
// and applies them by recreating the container.
func (s *containerService) SaveEnv(ctx context.Context, id uuid.UUID, env types.EnvVariables) error {
	for _, e := range env {
		err := s.vars.UpdateContainerVariableByID(ctx, e)
		if err != nil {
			return err
		}
	}
	return s.RecreateContainer(ctx, id)
}

func (s *containerService) GetContainerPorts(ctx context.Context, id uuid.UUID) (types.Ports, error) {
	return s.ports.GetContainerPorts(ctx, id)
}

func (s *containerService) PatchPort(ctx context.Context, p types.Port) error {
	err := p.Validate()
	if err != nil {
		return err
	}
	return s.ports.UpdateContainerPortByID(ctx, p)
}

func (s *containerService) DeletePort(ctx context.Context, id uuid.UUID) error {
	return s.ports.DeletePort(ctx, id)
}

func (s *containerService) CreatePort(ctx context.Context, p types.Port) error {
	p.ID = uuid.New()
	err := p.Validate()
	if err != nil {
		return err
	}
	return s.ports.CreatePort(ctx, p)
}

func (s *containerService) GetAllVersions(ctx context.Context, id uuid.UUID, useCache bool) ([]string, error) {
	c, err := s.containers.GetContainer(ctx, id)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.cacheImageTags[c.Image]

	if !useCache || !ok {
		versions, err := s.runner.GetAllVersions(ctx, *c)
		if err != nil {
			return nil, err
		}
		s.cacheImageTags[c.Image] = versions
	}

	return s.cacheImageTags[c.Image], nil
}

func (s *containerService) GetContainerInfo(ctx context.Context, id uuid.UUID) (map[string]any, error) {
	c, err := s.containers.GetContainer(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.runner.Info(ctx, *c)
}

func (s *containerService) WaitStatus(ctx context.Context, id uuid.UUID, status string) error {
	statusChan := make(chan string)
	defer close(statusChan)

	c, err := s.containers.GetContainer(ctx, id)
	if err != nil {
		return err
	}

	if c.Status == status {
		return nil
	}

	l := event.NewTempListener(func(e event.Event) error {
		switch e := e.(type) {
		case types.EventContainerStatusChange:
			if e.ContainerID != id {
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

func (s *containerService) GetLatestLogs(id uuid.UUID) ([]types.LogLine, error) {
	return s.logs.LoadBuffer(id)
}

func (s *containerService) GetTemplateByID(ctx context.Context, id string) (*types.Template, error) {
	serv, err := s.templates.Get(id)
	if err != nil {
		return nil, err
	}
	return &serv, nil
}

func (s *containerService) GetTemplates(ctx context.Context) []types.Template {
	return s.templates.GetAll()
}

func (s *containerService) setStatus(c *types.Container, status string) {
	if c.Status == status {
		return
	}

	// TODO: Status should be saved only in the memory.

	err := s.containers.SetStatus(context.Background(), c.ID, status)
	if err != nil {
		log.Error(err)
	}

	c.Status = status
	s.ctx.DispatchEvent(types.EventContainersChange{})
	s.ctx.DispatchEvent(types.EventContainerStatusChange{
		ContainerID: c.ID,
		Container:   *c,
		Name:        c.Name,
		Status:      status,
	})
}
