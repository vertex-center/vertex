package adapter

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/carlmjohnson/requests"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/pkg/vdocker"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
	"github.com/vertex-center/vlog"
)

type RunnerDockerAdapter struct{}

func NewRunnerDockerAdapter() RunnerDockerAdapter {
	return RunnerDockerAdapter{}
}

func (a RunnerDockerAdapter) Delete(instance *types.Instance) error {
	id, err := a.getContainerID(*instance)
	if err != nil {
		return err
	}

	apiError := api.Error{}
	err = requests.URL(config.Current.HostKernel).
		Pathf("/api/docker/container/%s", id).
		Delete().
		ErrorJSON(&apiError).
		Fetch(context.Background())

	if apiError.Code == api.ErrContainerNotFound {
		return ErrContainerNotFound
	}
	return err
}

func (a RunnerDockerAdapter) Start(instance *types.Instance, setStatus func(status string)) (io.ReadCloser, io.ReadCloser, error) {
	rErr, wErr := io.Pipe()
	rOut, wOut := io.Pipe()

	go func() {
		imageName := instance.DockerImageName()

		setStatus(types.InstanceStatusBuilding)

		instancePath := a.getPath(*instance)
		service := instance.Service

		// Build
		var err error
		var stdout, stderr io.ReadCloser
		if service.Methods.Docker.Dockerfile != nil {
			stdout, err = a.buildImageFromDockerfile(instancePath, imageName)
		} else if service.Methods.Docker.Image != nil {
			stdout, err = a.buildImageFromName(*service.Methods.Docker.Image)
		} else {
			err = errors.New("no Docker methods found")
		}
		if err != nil {
			return
		}

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer stdout.Close()

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				if scanner.Err() != nil {
					log.Error(scanner.Err())
					return
				}
				_, err := fmt.Fprintln(wOut, scanner.Text())
				if err != nil {
					return
				}
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

		log.Info("waiting for image to be built", vlog.String("uuid", instance.UUID.String()))

		wg.Wait()

		log.Info("image built", vlog.String("uuid", instance.UUID.String()))

		// Create
		id, err := a.getContainerID(*instance)
		if errors.Is(err, ErrContainerNotFound) {
			containerName := instance.DockerContainerName()

			log.Info("container doesn't exists, create it.",
				vlog.String("container_name", containerName),
			)

			options := types.CreateContainerOptions{
				ImageName:     imageName,
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
						if e.Type == "port" && e.Default == out {
							out = instance.Env[e.Name]
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
						source, err = filepath.Abs(path.Join(instancePath, "volumes", source))
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
					value := instance.Env[out]
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

			if service.Methods.Docker.Dockerfile != nil {
				id, err = a.createContainer(options)
			} else if service.Methods.Docker.Image != nil {
				options.ImageName = *service.Methods.Docker.Image
				id, err = a.createContainer(options)
			}
			if err != nil {
				return
			}
		} else if err != nil {
			return
		}

		// Start
		err = requests.URL(config.Current.HostKernel).
			Pathf("/api/docker/container/%s/start", id).
			Post().
			Fetch(context.Background())
		if err != nil {
			setStatus(types.InstanceStatusError)
			return
		}
		setStatus(types.InstanceStatusRunning)

		stdout, stderr, err = a.readLogs(id)
		if err != nil {
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

		a.watchForStatusChange(id, instance, setStatus)
	}()

	return rOut, rErr, nil
}

func (a RunnerDockerAdapter) Stop(instance *types.Instance) error {
	id, err := a.getContainerID(*instance)
	if err != nil {
		return err
	}

	return requests.URL(config.Current.HostKernel).
		Pathf("/api/docker/container/%s/stop", id).
		Post().
		Fetch(context.Background())
}

func (a RunnerDockerAdapter) Info(instance types.Instance) (map[string]any, error) {
	id, err := a.getContainerID(instance)
	if err != nil {
		return nil, err
	}

	var info types.InfoContainerResponse
	err = requests.URL(config.Current.HostKernel).
		Pathf("/api/docker/container/%s/info", id).
		ToJSON(&info).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	var imageInfo types.InfoImageResponse
	err = requests.URL(config.Current.HostKernel).
		Pathf("/api/docker/image/%s/info", info.Image).
		ToJSON(&imageInfo).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"container": info,
		"image":     imageInfo,
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

	var imageInfo types.InfoImageResponse
	err = requests.URL(config.Current.HostKernel).
		Pathf("/api/docker/%s/info", imageName).
		ToJSON(&imageInfo).
		Fetch(context.Background())
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

func (a RunnerDockerAdapter) getContainer(instance types.Instance) (types.Container, error) {
	var containers []types.Container
	err := requests.URL(config.Current.HostKernel).
		Path("/api/docker/containers").
		ToJSON(&containers).
		Fetch(context.Background())
	if err != nil {
		return types.Container{}, err
	}

	var dockerContainer *types.Container
	for _, c := range containers {
		name := c.Names[0]
		if name == "/"+instance.DockerContainerName() {
			dockerContainer = &c
			break
		}
	}

	if dockerContainer == nil {
		return types.Container{}, ErrContainerNotFound
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
	options := types.PullImageOptions{Image: imageName}

	req, err := requests.URL(config.Current.HostKernel).
		Path("/api/docker/image/pull").
		Post().
		BodyJSON(options).
		Request(context.Background())
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (a RunnerDockerAdapter) buildImageFromName(imageName string) (io.ReadCloser, error) {
	res, err := a.pullImage(imageName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a RunnerDockerAdapter) buildImageFromDockerfile(instancePath string, imageName string) (io.ReadCloser, error) {
	options := types.BuildImageOptions{
		Dir:        instancePath,
		Name:       imageName,
		Dockerfile: "Dockerfile",
	}

	req, err := requests.URL(config.Current.HostKernel).
		Pathf("/api/docker/image/build").
		Post().
		BodyJSON(options).
		Request(context.Background())
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
	}
	return res.Body, nil
}

func (a RunnerDockerAdapter) createContainer(options types.CreateContainerOptions) (string, error) {
	var res types.CreateContainerResponse
	err := requests.URL(config.Current.HostKernel).
		Pathf("/api/docker/container").
		Post().
		BodyJSON(options).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}

	for _, warn := range res.Warnings {
		log.Warn("warning while creating container",
			vlog.String("warning", warn),
		)
	}
	return res.ID, err
}

func (a RunnerDockerAdapter) watchForStatusChange(containerID string, instance *types.Instance, setStatus func(status string)) {
	go func() {
		err := requests.URL(config.Current.HostKernel).
			Pathf("/api/docker/container/%s/wait/%s", containerID, container.WaitConditionNotRunning).
			Fetch(context.Background())

		if err != nil {
			log.Error(err,
				vlog.String("uuid", instance.UUID.String()),
			)
			return
		}

		setStatus(types.InstanceStatusOff)
	}()
}

func (a RunnerDockerAdapter) readLogs(containerID string) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	var reqStdout, reqStderr *http.Request
	reqStdout, err = requests.URL(config.Current.HostKernel).
		Pathf("/api/docker/container/%s/logs/stdout", containerID).
		Request(context.Background())
	if err != nil {
		return
	}

	reqStderr, err = requests.URL(config.Current.HostKernel).
		Pathf("/api/docker/container/%s/logs/stderr", containerID).
		Request(context.Background())
	if err != nil {
		return
	}

	rOut, wOut := io.Pipe()
	rErr, wErr := io.Pipe()

	go func() {
		res, err := http.DefaultClient.Do(reqStdout)
		if err != nil {
			return
		}
		defer res.Body.Close()

		_, err = io.Copy(wOut, res.Body)
		if err != nil {
			return
		}
	}()

	go func() {
		res, err := http.DefaultClient.Do(reqStderr)
		if err != nil {
			return
		}
		defer res.Body.Close()

		_, err = io.Copy(wErr, res.Body)
		if err != nil {
			return
		}
	}()

	return rOut, rErr, nil
}

func (a RunnerDockerAdapter) getPath(instance types.Instance) string {
	base := storage.Path

	// If Vertex is running itself inside Docker, the instances are stored in the Vertex container volume.
	if vdocker.RunningInDocker() {
		var containers []types.Container
		err := requests.URL(config.Current.HostKernel).
			Path("/api/docker/containers").
			ToJSON(&containers).
			Fetch(context.Background())
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
