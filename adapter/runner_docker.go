package adapter

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"path"
	"path/filepath"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/pkg/vdocker"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type RunnerDockerAdapter struct {
	cli *client.Client
}

type dockerMessage struct {
	Stream string `json:"stream"`
}

func NewRunnerDockerAdapter() RunnerDockerAdapter {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Warn("couldn't connect with the Docker cli.",
			vlog.String("error", err.Error()),
		)

		return RunnerDockerAdapter{}
	}

	return RunnerDockerAdapter{
		cli: cli,
	}
}

func (a RunnerDockerAdapter) Delete(instance *types.Instance) error {
	id, err := a.getContainerID(*instance)
	if err != nil {
		return err
	}

	return a.cli.ContainerRemove(context.Background(), id, dockertypes.ContainerRemoveOptions{})
}

func (a RunnerDockerAdapter) Start(instance *types.Instance, onLog func(msg string), onErr func(msg string), setStatus func(status string)) error {
	imageName := instance.DockerImageName()

	setStatus(types.InstanceStatusBuilding)

	instancePath := a.getPath(*instance)
	service := instance.Service

	// Build
	var err error
	if service.Methods.Docker.Dockerfile != nil {
		err = a.buildImageFromDockerfile(instancePath, imageName, onLog)
	} else if service.Methods.Docker.Image != nil {
		err = a.buildImageFromName(*service.Methods.Docker.Image, onLog)
	} else {
		err = errors.New("no Docker methods found")
	}

	if err != nil {
		onErr(err.Error())
		return err
	}

	// Create
	id, err := a.getContainerID(*instance)
	if errors.Is(err, ErrContainerNotFound) {
		containerName := instance.DockerContainerName()

		log.Info("container doesn't exists, create it.",
			vlog.String("container_name", containerName),
		)

		options := createContainerOptions{
			imageName:     imageName,
			containerName: containerName,
			exposedPorts:  nat.PortSet{},
			portBindings:  nat.PortMap{},
			binds:         []string{},
			env:           []string{},
			capAdd:        []string{},
		}

		// exposedPorts and portBindings
		if service.Methods.Docker.Ports != nil {
			var all []string

			for in, out := range *service.Methods.Docker.Ports {
				for _, e := range service.Env {
					if e.Type == "port" && e.Default == out {
						out = instance.Env[e.Name]
						all = append(all, out+":"+in)
						break
					}
				}
			}

			var err error
			options.exposedPorts, options.portBindings, err = nat.ParsePortSpecs(all)
			if err != nil {
				return err
			}
		}

		// binds
		if service.Methods.Docker.Volumes != nil {
			for source, target := range *service.Methods.Docker.Volumes {
				if !strings.HasPrefix(source, "/") {
					source, err = filepath.Abs(path.Join(instancePath, "volumes", source))
				}
				if err != nil {
					return err
				}
				options.binds = append(options.binds, source+":"+target)
			}
		}

		// env
		if service.Methods.Docker.Environment != nil {
			for in, out := range *service.Methods.Docker.Environment {
				value := instance.Env[out]
				options.env = append(options.env, in+"="+value)
			}
		}

		// capAdd
		if service.Methods.Docker.Capabilities != nil {
			options.capAdd = *service.Methods.Docker.Capabilities
		}

		// sysctls
		if service.Methods.Docker.Sysctls != nil {
			options.sysctls = *service.Methods.Docker.Sysctls
		}

		if service.Methods.Docker.Dockerfile != nil {
			id, err = a.createContainer(options)
		} else if service.Methods.Docker.Image != nil {
			options.imageName = *service.Methods.Docker.Image
			id, err = a.createContainer(options)
		}
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Start
	err = a.cli.ContainerStart(context.Background(), id, dockertypes.ContainerStartOptions{})
	if err != nil {
		setStatus(types.InstanceStatusError)
		return err
	}
	setStatus(types.InstanceStatusRunning)

	a.watchForLogs(id, instance, onLog)
	a.watchForStatusChange(id, instance, setStatus)

	return nil
}

func (a RunnerDockerAdapter) Stop(instance *types.Instance) error {
	id, err := a.getContainerID(*instance)
	if err != nil {
		return err
	}

	return a.cli.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (a RunnerDockerAdapter) Info(instance types.Instance) (map[string]any, error) {
	id, err := a.getContainerID(instance)
	if err != nil {
		return nil, err
	}

	info, err := a.cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, err
	}

	imageInfo, _, err := a.cli.ImageInspectWithRaw(context.Background(), info.Image)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"container": map[string]any{
			"id":       info.ID,
			"name":     info.Name,
			"platform": info.Platform,
		},
		"image": map[string]any{
			"id":           imageInfo.ID,
			"architecture": imageInfo.Architecture,
			"os":           imageInfo.Os,
			"size":         imageInfo.Size,
			"virtual_size": imageInfo.VirtualSize,
			"tags":         imageInfo.RepoTags,
		},
	}, nil
}

func (a RunnerDockerAdapter) CheckForUpdates(instance *types.Instance) error {
	service := instance.Service

	if service.Methods.Docker.Image == nil {
		// TODO: Support Dockerfile updates
		return nil
	}

	imageName := *service.Methods.Docker.Image

	res, err := a.pullImage(imageName)
	if err != nil {
		return err
	}
	defer res.Close()

	imageInfo, _, err := a.cli.ImageInspectWithRaw(context.Background(), imageName)
	if err != nil {
		return err
	}

	latestImageID := imageInfo.ID

	currentImageID, err := a.getImageID(*instance)
	if err != nil {
		return err
	}

	if latestImageID == currentImageID {
		log.Info("already up-to-date",
			vlog.String("uuid", instance.UUID.String()),
		)
		instance.Update = nil
	} else {
		log.Info("a new update is available",
			vlog.String("uuid", instance.UUID.String()),
		)
		instance.Update = &types.InstanceUpdate{
			CurrentVersion: currentImageID,
			LatestVersion:  latestImageID,
		}
	}

	return nil
}

func (a RunnerDockerAdapter) HasUpdateAvailable(instance types.Instance) (bool, error) {
	//TODO implement me
	return false, nil
}

func (a RunnerDockerAdapter) getContainer(instance types.Instance) (dockertypes.Container, error) {
	containers, err := a.cli.ContainerList(context.Background(), dockertypes.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return dockertypes.Container{}, err
	}

	var dockerContainer *dockertypes.Container

	for _, c := range containers {
		name := c.Names[0]
		if name == "/"+instance.DockerContainerName() {
			dockerContainer = &c
			break
		}
	}

	if dockerContainer == nil {
		return dockertypes.Container{}, ErrContainerNotFound
	}

	return *dockerContainer, nil
}

func (a RunnerDockerAdapter) getContainerID(instance types.Instance) (string, error) {
	c, err := a.getContainer(instance)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func (a RunnerDockerAdapter) getImageID(instance types.Instance) (string, error) {
	c, err := a.getContainer(instance)
	if err != nil {
		return "", err
	}
	return c.ImageID, nil
}

func (a RunnerDockerAdapter) pullImage(imageName string) (io.ReadCloser, error) {
	res, err := a.cli.ImagePull(context.Background(), imageName, dockertypes.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a RunnerDockerAdapter) buildImageFromName(imageName string, onMsg func(msg string)) error {
	res, err := a.pullImage(imageName)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(res)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return scanner.Err()
		}
		onMsg(scanner.Text())
	}

	return nil
}

func (a RunnerDockerAdapter) buildImageFromDockerfile(instancePath string, imageName string, onMsg func(msg string)) error {
	buildOptions := dockertypes.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Remove:     true,
	}

	reader, err := archive.TarWithOptions(instancePath, &archive.TarOptions{
		ExcludePatterns: []string{".git/**/*"},
	})
	if err != nil {
		return err
	}

	res, err := a.cli.ImageBuild(context.Background(), reader, buildOptions)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return scanner.Err()
		}
		msg := dockerMessage{}
		err := json.Unmarshal(scanner.Bytes(), &msg)
		if err != nil {
			log.Warn("Failed to parse message",
				vlog.String("message", scanner.Text()),
			)
		} else {
			if msg.Stream != "" {
				onMsg(msg.Stream)
			}
		}
	}

	log.Info("Docker build: success.")
	return nil
}

type createContainerOptions struct {
	imageName     string
	containerName string
	exposedPorts  nat.PortSet
	portBindings  nat.PortMap
	binds         []string
	env           []string
	capAdd        []string
	sysctls       map[string]string
}

func (a RunnerDockerAdapter) createContainer(options createContainerOptions) (string, error) {
	config := container.Config{
		Image:        options.imageName,
		ExposedPorts: options.exposedPorts,
		Env:          options.env,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}

	hostConfig := container.HostConfig{
		Binds:        options.binds,
		PortBindings: options.portBindings,
		CapAdd:       options.capAdd,
		Sysctls:      options.sysctls,
	}

	res, err := a.cli.ContainerCreate(context.Background(), &config, &hostConfig, nil, nil, options.containerName)
	for _, warn := range res.Warnings {
		log.Warn("warning while creating container",
			vlog.String("warning", warn),
		)
	}
	return res.ID, err
}

func (a RunnerDockerAdapter) watchForStatusChange(containerID string, instance *types.Instance, setStatus func(status string)) {
	go func() {
		resChan, errChan := a.cli.ContainerWait(context.Background(), containerID, container.WaitConditionNotRunning)

		select {
		case err := <-errChan:
			if err != nil {
				log.Error(err,
					vlog.String("uuid", instance.UUID.String()),
				)
			}
		case status := <-resChan:
			log.Info("container exited",
				vlog.String("uuid", instance.UUID.String()),
				vlog.Int64("status", status.StatusCode),
			)
			setStatus(types.InstanceStatusOff)
		}
	}()
}

func (a RunnerDockerAdapter) watchForLogs(containerID string, instance *types.Instance, onLog func(msg string)) {
	go func() {
		logs, err := a.cli.ContainerLogs(context.Background(), containerID, dockertypes.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: false,
			Follow:     true,
			Tail:       "0",
		})
		if err != nil {
			log.Error(err,
				vlog.String("uuid", instance.UUID.String()),
			)
		}

		scanner := bufio.NewScanner(logs)
		for scanner.Scan() {
			onLog(scanner.Text())
		}
		_ = logs.Close()
		log.Info("logs pipe closed",
			vlog.String("uuid", instance.UUID.String()),
		)
	}()
}

func (a RunnerDockerAdapter) getPath(instance types.Instance) string {
	base := storage.Path

	// If Vertex is running itself inside Docker, the instances are stored in the Vertex container volume.
	if vdocker.RunningInDocker() {
		containers, err := a.cli.ContainerList(context.Background(), dockertypes.ContainerListOptions{
			All: true,
		})
		if err != nil {
			log.Error(err)
		} else {
			for _, c := range containers {
				// find the docker container that has a volume /live, which is the Vertex container.
				for _, m := range c.Mounts {
					if m.Destination == "/live" {
						base = m.Source
					}
				}
			}
		}
	}

	return path.Join(base, "instances", instance.UUID.String())
}
