package adapter

import (
	"bufio"
	"context"
	"encoding/json"
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
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/google/go-containerregistry/pkg/crane"
	containerstypes "github.com/vertex-center/vertex/apps/containers/types"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/pkg/vdocker"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
	"github.com/vertex-center/vlog"
)

type ContainerRunnerDockerAdapter struct{}

func NewContainerRunnerFSAdapter() ContainerRunnerDockerAdapter {
	return ContainerRunnerDockerAdapter{}
}

func (a ContainerRunnerDockerAdapter) Delete(inst *containerstypes.Container) error {
	id, err := a.getContainerID(*inst)
	if err != nil {
		return err
	}

	apiError := router.Error{}
	err = requests.URL(config.Current.KernelURL()).
		Pathf("/api/docker/container/%s", id).
		Delete().
		ErrorJSON(&apiError).
		Fetch(context.Background())

	if apiError.Code == api.ErrContainerNotFound {
		return ErrContainerNotFound
	}
	return err
}

func (a ContainerRunnerDockerAdapter) Start(inst *containerstypes.Container, setStatus func(status string)) (io.ReadCloser, io.ReadCloser, error) {
	rErr, wErr := io.Pipe()
	rOut, wOut := io.Pipe()

	go func() {
		imageName := inst.DockerImageVertexName()

		setStatus(containerstypes.ContainerStatusBuilding)

		containerPath := a.getPath(*inst)
		service := inst.Service

		log.Debug("building image", vlog.String("image", imageName))

		// Build
		var err error
		var stdout, stderr io.ReadCloser
		if service.Methods.Docker.Dockerfile != nil {
			stdout, err = a.buildImageFromDockerfile(containerPath, imageName)
		} else if service.Methods.Docker.Image != nil {
			stdout, err = a.buildImageFromName(inst.GetImageNameWithTag())
		} else {
			err = errors.New("no Docker methods found")
		}
		if err != nil {
			log.Error(err)
			setStatus(containerstypes.ContainerStatusError)
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

				progress := containerstypes.DownloadProgress{
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
					setStatus(containerstypes.ContainerStatusError)
					return
				}
			}
			if scanner.Err() != nil {
				log.Error(scanner.Err(),
					vlog.String("uuid", inst.UUID.String()))
				setStatus(containerstypes.ContainerStatusError)
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
						if e.Type == "port" && e.Default == out {
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
						source, err = filepath.Abs(path.Join(containerPath, "volumes", source))
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
		err = requests.URL(config.Current.KernelURL()).
			Pathf("/api/docker/container/%s/start", id).
			Post().
			Fetch(context.Background())
		if err != nil {
			setStatus(containerstypes.ContainerStatusError)
			return
		}
		setStatus(containerstypes.ContainerStatusRunning)

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

		err = a.WaitCondition(inst, types.WaitContainerCondition(container.WaitConditionNotRunning))
		if err != nil {
			log.Error(err)
			setStatus(containerstypes.ContainerStatusError)
		} else {
			setStatus(containerstypes.ContainerStatusOff)
		}
	}()

	return rOut, rErr, nil
}

func (a ContainerRunnerDockerAdapter) Stop(inst *containerstypes.Container) error {
	id, err := a.getContainerID(*inst)
	if err != nil {
		return err
	}

	return requests.URL(config.Current.KernelURL()).
		Pathf("/api/docker/container/%s/stop", id).
		Post().
		Fetch(context.Background())
}

func (a ContainerRunnerDockerAdapter) Info(inst containerstypes.Container) (map[string]any, error) {
	id, err := a.getContainerID(inst)
	if err != nil {
		return nil, err
	}

	var info types.InfoContainerResponse
	err = requests.URL(config.Current.KernelURL()).
		Pathf("/api/docker/container/%s/info", id).
		ToJSON(&info).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	var imageInfo types.InfoImageResponse
	err = requests.URL(config.Current.KernelURL()).
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

func (a ContainerRunnerDockerAdapter) CheckForUpdates(inst *containerstypes.Container) error {
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

	var imageInfo types.InfoImageResponse
	err = requests.URL(config.Current.KernelURL()).
		Pathf("/api/docker/%s/info", imageName).
		ToJSON(&imageInfo).
		Fetch(context.Background())
	if err != nil {
		return err
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
		inst.Update = &containerstypes.ContainerUpdate{
			CurrentVersion: currentImageID,
			LatestVersion:  latestImageID,
		}
	}

	return nil
}

func (a ContainerRunnerDockerAdapter) GetAllVersions(inst containerstypes.Container) ([]string, error) {
	if inst.Service.Methods.Docker == nil {
		return nil, errors.New("no Docker methods found")
	}
	image := *inst.Service.Methods.Docker.Image
	log.Debug("querying all versions of image",
		vlog.String("image", image),
	)
	return crane.ListTags(image)
}

func (a ContainerRunnerDockerAdapter) HasUpdateAvailable(inst containerstypes.Container) (bool, error) {
	//TODO implement me
	return false, nil
}

func (a ContainerRunnerDockerAdapter) WaitCondition(inst *containerstypes.Container, cond types.WaitContainerCondition) error {
	id, err := a.getContainerID(*inst)
	if err != nil {
		return err
	}

	return requests.URL(config.Current.KernelURL()).
		Pathf("/api/docker/container/%s/wait/%s", id, cond).
		Fetch(context.Background())
}

func (a ContainerRunnerDockerAdapter) getContainer(inst containerstypes.Container) (types.Container, error) {
	var containers []types.Container
	err := requests.URL(config.Current.KernelURL()).
		Path("/api/docker/containers").
		ToJSON(&containers).
		Fetch(context.Background())
	if err != nil {
		return types.Container{}, err
	}

	var dockerContainer *types.Container
	for _, c := range containers {
		name := c.Names[0]
		if name == "/"+inst.DockerContainerName() {
			dockerContainer = &c
			break
		}
	}

	if dockerContainer == nil {
		return types.Container{}, ErrContainerNotFound
	}

	return *dockerContainer, nil
}

func (a ContainerRunnerDockerAdapter) getContainerID(inst containerstypes.Container) (string, error) {
	c, err := a.getContainer(inst)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func (a ContainerRunnerDockerAdapter) getImageID(inst containerstypes.Container) (string, error) {
	c, err := a.getContainer(inst)
	if err != nil {
		return "", err
	}
	return c.ImageID, nil
}

func (a ContainerRunnerDockerAdapter) pullImage(imageName string) (io.ReadCloser, error) {
	options := types.PullImageOptions{Image: imageName}

	req, err := requests.URL(config.Current.KernelURL()).
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
	} else if res.StatusCode >= 200 && res.StatusCode < 300 {
		return res.Body, nil
	}
	return nil, errors.New("failed to pull image")
}

func (a ContainerRunnerDockerAdapter) buildImageFromName(imageName string) (io.ReadCloser, error) {
	res, err := a.pullImage(imageName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a ContainerRunnerDockerAdapter) buildImageFromDockerfile(containerPath string, imageName string) (io.ReadCloser, error) {
	options := types.BuildImageOptions{
		Dir:        containerPath,
		Name:       imageName,
		Dockerfile: "Dockerfile",
	}

	req, err := requests.URL(config.Current.KernelURL()).
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

func (a ContainerRunnerDockerAdapter) createContainer(options types.CreateContainerOptions) (string, error) {
	var res types.CreateContainerResponse
	err := requests.URL(config.Current.KernelURL()).
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

func (a ContainerRunnerDockerAdapter) readLogs(containerID string) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	var reqStdout, reqStderr *http.Request
	reqStdout, err = requests.URL(config.Current.KernelURL()).
		Pathf("/api/docker/container/%s/logs/stdout", containerID).
		Request(context.Background())
	if err != nil {
		return
	}

	reqStderr, err = requests.URL(config.Current.KernelURL()).
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

func (a ContainerRunnerDockerAdapter) getPath(inst containerstypes.Container) string {
	base := storage.Path

	// If Vertex is running itself inside Docker, the containers are stored in the Vertex container volume.
	if vdocker.RunningInDocker() {
		var containers []types.Container
		err := requests.URL(config.Current.KernelURL()).
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

	return path.Join(base, "apps", "vx-containers", inst.UUID.String())
}
