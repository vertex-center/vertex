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
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/pkg/vdocker"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
	"github.com/vertex-center/vlog"
)

type InstanceRunnerDockerAdapter struct{}

func NewInstanceRunnerFSAdapter() InstanceRunnerDockerAdapter {
	return InstanceRunnerDockerAdapter{}
}

func (a InstanceRunnerDockerAdapter) Delete(inst *instancestypes.Instance) error {
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

func (a InstanceRunnerDockerAdapter) Start(inst *instancestypes.Instance, setStatus func(status string)) (io.ReadCloser, io.ReadCloser, error) {
	rErr, wErr := io.Pipe()
	rOut, wOut := io.Pipe()

	go func() {
		imageName := inst.DockerImageVertexName()

		setStatus(instancestypes.InstanceStatusBuilding)

		instancePath := a.getPath(*inst)
		service := inst.Service

		// Build
		var err error
		var stdout, stderr io.ReadCloser
		if service.Methods.Docker.Dockerfile != nil {
			stdout, err = a.buildImageFromDockerfile(instancePath, imageName)
		} else if service.Methods.Docker.Image != nil {
			stdout, err = a.buildImageFromName(inst.GetImageNameWithTag())
		} else {
			err = errors.New("no Docker methods found")
		}
		if err != nil {
			log.Error(err)
			setStatus(instancestypes.InstanceStatusError)
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

				var msg jsonmessage.JSONMessage
				err := json.Unmarshal(scanner.Bytes(), &msg)
				if err != nil {
					log.Error(err)
					continue
				}

				progress := instancestypes.DownloadProgress{
					ID:     msg.ID,
					Status: msg.Status,
				}

				if msg.Progress != nil {
					progress.Current = msg.Progress.Current
					progress.Total = msg.Progress.Total
				}

				progressJSON, err := json.Marshal(progress)
				if err != nil {
					log.Error(err)
					continue
				}

				_, err = fmt.Fprintf(wOut, "%s %s\n", "DOWNLOAD", progressJSON)
				if err != nil {
					log.Error(err)
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
			setStatus(instancestypes.InstanceStatusError)
			return
		}
		setStatus(instancestypes.InstanceStatusRunning)

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

		a.WaitStatus(id, inst, container.WaitConditionNotRunning)
		setStatus(instancestypes.InstanceStatusOff)
	}()

	return rOut, rErr, nil
}

func (a InstanceRunnerDockerAdapter) Stop(inst *instancestypes.Instance) error {
	id, err := a.getContainerID(*inst)
	if err != nil {
		return err
	}

	return requests.URL(config.Current.KernelURL()).
		Pathf("/api/docker/container/%s/stop", id).
		Post().
		Fetch(context.Background())
}

func (a InstanceRunnerDockerAdapter) Info(inst instancestypes.Instance) (map[string]any, error) {
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

func (a InstanceRunnerDockerAdapter) CheckForUpdates(inst *instancestypes.Instance) error {
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
		inst.Update = &instancestypes.InstanceUpdate{
			CurrentVersion: currentImageID,
			LatestVersion:  latestImageID,
		}
	}

	return nil
}

func (a InstanceRunnerDockerAdapter) GetAllVersions(inst instancestypes.Instance) ([]string, error) {
	if inst.Service.Methods.Docker == nil {
		return nil, errors.New("no Docker methods found")
	}
	image := *inst.Service.Methods.Docker.Image
	log.Debug("querying all versions of image",
		vlog.String("image", image),
	)
	return crane.ListTags(image)
}

func (a InstanceRunnerDockerAdapter) HasUpdateAvailable(inst instancestypes.Instance) (bool, error) {
	//TODO implement me
	return false, nil
}

func (a InstanceRunnerDockerAdapter) getContainer(inst instancestypes.Instance) (types.Container, error) {
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

func (a InstanceRunnerDockerAdapter) getContainerID(inst instancestypes.Instance) (string, error) {
	c, err := a.getContainer(inst)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func (a InstanceRunnerDockerAdapter) getImageID(inst instancestypes.Instance) (string, error) {
	c, err := a.getContainer(inst)
	if err != nil {
		return "", err
	}
	return c.ImageID, nil
}

func (a InstanceRunnerDockerAdapter) pullImage(imageName string) (io.ReadCloser, error) {
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

func (a InstanceRunnerDockerAdapter) buildImageFromName(imageName string) (io.ReadCloser, error) {
	res, err := a.pullImage(imageName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a InstanceRunnerDockerAdapter) buildImageFromDockerfile(instancePath string, imageName string) (io.ReadCloser, error) {
	options := types.BuildImageOptions{
		Dir:        instancePath,
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

func (a InstanceRunnerDockerAdapter) createContainer(options types.CreateContainerOptions) (string, error) {
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

func (a InstanceRunnerDockerAdapter) WaitStatus(containerID string, inst *instancestypes.Instance, condition container.WaitCondition) {
	go func() {
		err := requests.URL(config.Current.KernelURL()).
			Pathf("/api/docker/container/%s/wait/%s", containerID, condition).
			Fetch(context.Background())

		if err != nil {
			log.Error(err,
				vlog.String("uuid", inst.UUID.String()),
			)
			return
		}
	}()
}

func (a InstanceRunnerDockerAdapter) readLogs(containerID string) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
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

func (a InstanceRunnerDockerAdapter) getPath(inst instancestypes.Instance) string {
	base := storage.Path

	// If Vertex is running itself inside Docker, the instances are stored in the Vertex container volume.
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

	return path.Join(base, "instances", inst.UUID.String())
}
