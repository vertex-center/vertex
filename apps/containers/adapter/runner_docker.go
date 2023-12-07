package adapter

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/uuid"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/vdocker"
	"github.com/vertex-center/vlog"
)

type containerRunnerDockerAdapter struct{}

func NewContainerRunnerFSAdapter() port.ContainerRunnerAdapter {
	return containerRunnerDockerAdapter{}
}

func (a containerRunnerDockerAdapter) DeleteContainer(inst *types.Container) error {
	id, err := a.getContainerID(*inst)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient()
	apiError := cli.DeleteContainer(context.Background(), id)
	if apiError != nil {
		return apiError.RouterError()
	}
	return nil
}

func (a containerRunnerDockerAdapter) DeleteMounts(inst *types.Container) error {
	id, err := a.getContainerID(*inst)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient()
	apiError := cli.DeleteMounts(context.Background(), id)
	if apiError != nil {
		return apiError.RouterError()
	}
	return nil
}

func (a containerRunnerDockerAdapter) Start(inst *types.Container, setStatus func(status string)) (io.ReadCloser, io.ReadCloser, error) {
	rErr, wErr := io.Pipe()
	rOut, wOut := io.Pipe()

	go func() {
		imageName := inst.DockerImageVertexName()

		setStatus(types.ContainerStatusBuilding)

		service := inst.Service

		log.Debug("building image", vlog.String("image", imageName))

		// Build
		var err error
		var stdout, stderr io.ReadCloser
		if service.Methods.Docker.Dockerfile != nil {
			containerPath := a.getContainerPath(inst.UUID)
			stdout, err = a.buildImageFromDockerfile(containerPath, imageName)
		} else if service.Methods.Docker.Image != nil {
			stdout, err = a.buildImageFromName(inst.GetImageNameWithTag())
		} else {
			err = errors.New("no Docker methods found")
		}
		if err != nil {
			log.Error(err)
			setStatus(types.ContainerStatusError)
			return
		}

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer stdout.Close()

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				var msg jsonmessage.JSONMessage
				err := json.Unmarshal(scanner.Bytes(), &msg)
				if err != nil {
					log.Error(err,
						vlog.String("text", scanner.Text()),
						vlog.String("uuid", inst.UUID.String()))
					continue
				}

				progress := types.DownloadProgress{
					ID:     msg.ID,
					Status: msg.Status,
				}

				if msg.Progress != nil {
					progress.Current = msg.Progress.Current
					progress.Total = msg.Progress.Total
				}

				progressJSON, err := json.Marshal(progress)
				if err != nil {
					log.Error(err,
						vlog.String("text", scanner.Text()),
						vlog.String("uuid", inst.UUID.String()))
					continue
				}

				_, err = fmt.Fprintf(wOut, "%s %s\n", "DOWNLOAD", progressJSON)
				if err != nil {
					log.Error(err,
						vlog.String("text", scanner.Text()),
						vlog.String("uuid", inst.UUID.String()))
					setStatus(types.ContainerStatusError)
					return
				}
			}
			if scanner.Err() != nil {
				log.Error(scanner.Err(),
					vlog.String("uuid", inst.UUID.String()))
				setStatus(types.ContainerStatusError)
				return
			}
		}()

		//wg.Add(1)
		//go func() {
		//	defer wg.Done()
		//	defer stderr.Close()
		//	_, err := io.Copy(wErr, stderr)
		//	if err != nil {
		//		log.Error(err)
		//		return
		//	}
		//}()

		log.Info("waiting for image to be built", vlog.String("uuid", inst.UUID.String()))

		wg.Wait()

		log.Info("image built", vlog.String("uuid", inst.UUID.String()))

		// Create
		id, err := a.getContainerID(*inst)
		if errors.Is(err, ErrContainerNotFound) {
			containerName := inst.DockerContainerName()

			log.Info("container doesn't exists, create it.",
				vlog.String("container_name", containerName),
			)

			options := types.CreateContainerOptions{
				ContainerName: containerName,
				ExposedPorts:  nat.PortSet{},
				PortBindings:  nat.PortMap{},
				Binds:         []string{},
				Env:           []string{},
				CapAdd:        []string{},
			}

			// exposedPorts and portBindings
			if service.Methods.Docker.Ports != nil {
				var all []string

				for in, out := range *service.Methods.Docker.Ports {
					for _, e := range service.Env {
						if e.Type == "port" && e.Name == out {
							out = inst.Env[e.Name]
							all = append(all, out+":"+in)
							break
						}
					}
				}

				options.ExposedPorts, options.PortBindings, err = nat.ParsePortSpecs(all)
				if err != nil {
					return
				}
			}

			// binds
			if service.Methods.Docker.Volumes != nil {
				for source, target := range *service.Methods.Docker.Volumes {
					if !strings.HasPrefix(source, "/") {
						volumePath := a.getVolumePath(inst.UUID)
						source, err = filepath.Abs(path.Join(volumePath, source))
					}
					if err != nil {
						return
					}
					options.Binds = append(options.Binds, source+":"+target)
				}
			}

			// env
			if service.Methods.Docker.Environment != nil {
				for in, out := range *service.Methods.Docker.Environment {
					value := inst.Env[out]
					options.Env = append(options.Env, in+"="+value)
				}
			}

			// capAdd
			if service.Methods.Docker.Capabilities != nil {
				options.CapAdd = *service.Methods.Docker.Capabilities
			}

			// sysctls
			if service.Methods.Docker.Sysctls != nil {
				options.Sysctls = *service.Methods.Docker.Sysctls
			}

			// cmd
			if service.Methods.Docker.Cmd != nil {
				options.Cmd = strings.Split(*service.Methods.Docker.Cmd, " ")
			}

			if service.Methods.Docker.Dockerfile != nil {
				options.ImageName = inst.DockerImageVertexName()
				id, err = a.createContainer(options)
			} else if service.Methods.Docker.Image != nil {
				options.ImageName = inst.GetImageNameWithTag()
				id, err = a.createContainer(options)
			}
			if err != nil {
				return
			}
		} else if err != nil {
			return
		}

		// Start
		cli := containersapi.NewContainersKernelClient()
		apiError := cli.StartContainer(context.Background(), id)
		if apiError != nil {
			log.Error(apiError.RouterError())
			setStatus(types.ContainerStatusError)
			return
		}
		setStatus(types.ContainerStatusRunning)

		stdout, stderr, err = a.readLogs(id)
		if err != nil {
			log.Error(err)
			return
		}

		go func() {
			defer stdout.Close()
			defer wOut.Close()

			_, err := io.Copy(wOut, stdout)
			if err != nil {
				log.Error(err)
				return
			}
		}()

		go func() {
			defer stderr.Close()
			defer wErr.Close()

			_, err := io.Copy(wOut, stdout)
			if err != nil {
				log.Error(err)
				return
			}
		}()

		err = a.WaitCondition(inst, types.WaitContainerCondition(container.WaitConditionNotRunning))
		if err != nil {
			log.Error(err)
			setStatus(types.ContainerStatusError)
		} else {
			setStatus(types.ContainerStatusOff)
		}
	}()

	return rOut, rErr, nil
}

func (a containerRunnerDockerAdapter) Stop(inst *types.Container) error {
	id, err := a.getContainerID(*inst)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient()
	apiError := cli.StopContainer(context.Background(), id)
	if apiError != nil {
		return apiError.RouterError()
	}
	return nil
}

func (a containerRunnerDockerAdapter) Info(inst types.Container) (map[string]any, error) {
	id, err := a.getContainerID(inst)
	if err != nil {
		return nil, err
	}

	cli := containersapi.NewContainersKernelClient()
	info, apiError := cli.GetContainerInfo(context.Background(), id)
	if apiError != nil {
		return nil, apiError.RouterError()
	}

	imageInfo, apiError := cli.GetImageInfo(context.Background(), info.Image)
	if apiError != nil {
		return nil, apiError.RouterError()
	}

	return map[string]any{
		"container": info,
		"image":     imageInfo,
	}, nil
}

func (a containerRunnerDockerAdapter) CheckForUpdates(inst *types.Container) error {
	service := inst.Service

	if service.Methods.Docker.Image == nil {
		// TODO: Support Dockerfile updates
		return nil
	}

	imageName := inst.GetImageNameWithTag()

	res, err := a.pullImage(imageName)
	if err != nil {
		return err
	}
	defer res.Close()

	client := containersapi.NewContainersKernelClient()
	imageInfo, apiError := client.GetImageInfo(context.Background(), imageName)
	if apiError != nil {
		return apiError.RouterError()
	}

	latestImageID := imageInfo.ID

	currentImageID, err := a.getImageID(*inst)
	if err != nil {
		return err
	}

	if latestImageID == currentImageID {
		log.Info("already up-to-date",
			vlog.String("uuid", inst.UUID.String()),
		)
		inst.Update = nil
	} else {
		log.Info("a new update is available",
			vlog.String("uuid", inst.UUID.String()),
		)
		inst.Update = &types.ContainerUpdate{
			CurrentVersion: currentImageID,
			LatestVersion:  latestImageID,
		}
	}

	return nil
}

func (a containerRunnerDockerAdapter) GetAllVersions(inst types.Container) ([]string, error) {
	if inst.Service.Methods.Docker == nil {
		return nil, errors.New("no Docker methods found")
	}
	image := *inst.Service.Methods.Docker.Image
	log.Debug("querying all versions of image",
		vlog.String("image", image),
	)
	return crane.ListTags(image)
}

func (a containerRunnerDockerAdapter) HasUpdateAvailable(inst types.Container) (bool, error) {
	//TODO implement me
	return false, nil
}

func (a containerRunnerDockerAdapter) WaitCondition(inst *types.Container, cond types.WaitContainerCondition) error {
	id, err := a.getContainerID(*inst)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient()
	apiError := cli.WaitContainer(context.Background(), id, string(cond))
	if apiError != nil {
		return apiError.RouterError()
	}
	return nil
}

func (a containerRunnerDockerAdapter) getContainer(inst types.Container) (types.DockerContainer, error) {
	cli := containersapi.NewContainersKernelClient()
	containers, apiError := cli.GetContainers(context.Background())
	if apiError != nil {
		return types.DockerContainer{}, apiError.RouterError()
	}

	var dockerContainer *types.DockerContainer
	for _, c := range containers {
		name := c.Names[0]
		if name == "/"+inst.DockerContainerName() {
			dockerContainer = &c
			break
		}
	}

	if dockerContainer == nil {
		return types.DockerContainer{}, ErrContainerNotFound
	}

	return *dockerContainer, nil
}

func (a containerRunnerDockerAdapter) getContainerID(inst types.Container) (string, error) {
	c, err := a.getContainer(inst)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func (a containerRunnerDockerAdapter) getImageID(inst types.Container) (string, error) {
	c, err := a.getContainer(inst)
	if err != nil {
		return "", err
	}
	return c.ImageID, nil
}

func (a containerRunnerDockerAdapter) pullImage(imageName string) (io.ReadCloser, error) {
	options := types.PullImageOptions{Image: imageName}

	cli := containersapi.NewContainersKernelClient()
	res, apiError := cli.PullImage(context.Background(), options)
	if apiError != nil {
		return nil, apiError.RouterError()
	}
	return res, nil
}

func (a containerRunnerDockerAdapter) buildImageFromName(imageName string) (io.ReadCloser, error) {
	res, err := a.pullImage(imageName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a containerRunnerDockerAdapter) buildImageFromDockerfile(containerPath string, imageName string) (io.ReadCloser, error) {
	options := types.BuildImageOptions{
		Dir:        containerPath,
		Name:       imageName,
		Dockerfile: "Dockerfile",
	}

	cli := containersapi.NewContainersKernelClient()
	res, apiError := cli.BuildImage(context.Background(), options)
	if apiError != nil {
		return nil, apiError.RouterError()
	}
	return res, nil
}

func (a containerRunnerDockerAdapter) createContainer(options types.CreateContainerOptions) (string, error) {
	cli := containersapi.NewContainersKernelClient()
	res, apiError := cli.CreateContainer(context.Background(), options)
	if apiError != nil {
		return "", apiError.RouterError()
	}

	for _, warn := range res.Warnings {
		log.Warn("warning while creating container",
			vlog.String("warning", warn),
		)
	}
	return res.ID, nil
}

func (a containerRunnerDockerAdapter) readLogs(containerID string) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	cli := containersapi.NewContainersKernelClient()

	stdout, apiError := cli.GetContainerStdout(context.Background(), containerID)
	if apiError != nil {
		return nil, nil, apiError.RouterError()
	}

	stderr, apiError = cli.GetContainerStderr(context.Background(), containerID)
	if apiError != nil {
		return nil, nil, apiError.RouterError()
	}

	return stdout, stderr, nil
}

func (a containerRunnerDockerAdapter) getVolumePath(uuid uuid.UUID) string {
	appPath := a.getAppPath("live_docker")
	return path.Join(appPath, "volumes", uuid.String())
}

func (a containerRunnerDockerAdapter) getContainerPath(uuid uuid.UUID) string {
	appPath := a.getAppPath("live")
	return path.Join(appPath, "containers", uuid.String())
}

func (a containerRunnerDockerAdapter) getAppPath(base string) string {
	// If Vertex is running itself inside Docker, the containers are stored in the Vertex container volume.
	if vdocker.RunningInDocker() {
		cli := containersapi.NewContainersKernelClient()
		containers, apiError := cli.GetContainers(context.Background())
		if apiError != nil {
			log.Error(apiError.RouterError())
		} else {
			for _, c := range containers {
				// find the docker container that has a volume /live, which is the Vertex container.
				for _, m := range c.Mounts {
					if m.Destination == "/"+base {
						base = m.Source
					}
				}
			}
		}
	}

	return path.Join(base, "apps", "containers")
}
