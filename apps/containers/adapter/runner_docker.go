package adapter

import (
	"bufio"
	"context"
	"encoding/json"
	goerrors "errors"
	"fmt"
	"io"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/juju/errors"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/containers/core/types/builder"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vlog"
)

type runnerDockerAdapter struct{}

func NewRunnerDockerAdapter() port.RunnerAdapter {
	return runnerDockerAdapter{}
}

func (a runnerDockerAdapter) DeleteContainer(ctx context.Context, c *types.Container, volumes []string) error {
	id, err := a.getContainerID(ctx, *c)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	err = cli.DeleteContainer(context.Background(), id)
	if err != nil {
		return err
	}

	var errs []error
	for _, v := range volumes {
		err = cli.DeleteVolume(context.Background(), v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return goerrors.Join(errs...)
}

func (a runnerDockerAdapter) DeleteMounts(ctx context.Context, c *types.Container) error {
	id, err := a.getContainerID(ctx, *c)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	return cli.DeleteMounts(context.Background(), id)
}

func (a runnerDockerAdapter) Start(
	ctx context.Context,
	c *types.Container,
	ports types.Ports,
	volumes types.Volumes,
	env types.EnvVariables,
	caps types.Capabilities,
	sysctls types.Sysctls,
	setStatus func(status string),
) (io.ReadCloser, io.ReadCloser, error) {
	rErr, wErr := io.Pipe()
	rOut, wOut := io.Pipe()

	go func() {
		imageName := c.DockerImageVertexName()

		setStatus(types.ContainerStatusBuilding)

		log.Debug("building image", vlog.String("image", imageName))

		// Build
		stdout, err := a.buildImageFromName(ctx, c.GetImageNameWithTag())
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
						vlog.String("id", c.ID.String()))
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
						vlog.String("id", c.ID.String()))
					continue
				}

				_, err = fmt.Fprintf(wOut, "%s %s\n", "DOWNLOAD", progressJSON)
				if err != nil {
					log.Error(err,
						vlog.String("text", scanner.Text()),
						vlog.String("id", c.ID.String()))
					setStatus(types.ContainerStatusError)
					return
				}
			}
			if scanner.Err() != nil {
				log.Error(scanner.Err(),
					vlog.String("id", c.ID.String()))
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

		log.Info("waiting for image to be built", vlog.String("uuid", c.ID.String()))

		wg.Wait()

		log.Info("image built", vlog.String("uuid", c.ID.String()))

		// Create
		id, err := a.getContainerID(ctx, *c)
		if errors.Is(err, errors.NotFound) {
			containerName := c.DockerContainerName()

			log.Info("container doesn't exists, create it.", vlog.String("name", containerName))

			opts := builder.NewContainerOpts().
				WithName(containerName).
				WithImage(c.GetImageNameWithTag()).
				WithEnv(env).
				WithCaps(caps).
				WithSysctls(sysctls).
				WithPorts(ports, env).
				WithVolumes(volumes).
				WithCommand(c.Command).
				Build()

			id, err = a.createContainer(ctx, opts)
			if err != nil {
				log.Error(err)
				return
			}
		} else if err != nil {
			log.Error(err)
			return
		}

		// Start
		cli := containersapi.NewContainersKernelClient(ctx)
		err = cli.StartContainer(context.Background(), id)
		if err != nil {
			log.Error(err)
			setStatus(types.ContainerStatusError)
			return
		}
		setStatus(types.ContainerStatusRunning)

		var stderr io.ReadCloser
		stdout, stderr, err = a.readLogs(ctx, id)
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

		err = a.WaitCondition(ctx, c, types.WaitContainerCondition(container.WaitConditionNotRunning))
		if err != nil {
			log.Error(err)
			setStatus(types.ContainerStatusError)
		} else {
			setStatus(types.ContainerStatusOff)
		}
	}()

	return rOut, rErr, nil
}

func (a runnerDockerAdapter) Stop(ctx context.Context, c *types.Container) error {
	id, err := a.getContainerID(ctx, *c)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	return cli.StopContainer(context.Background(), id)
}

func (a runnerDockerAdapter) Info(ctx context.Context, c types.Container) (map[string]any, error) {
	id, err := a.getContainerID(ctx, c)
	if err != nil {
		return nil, err
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	info, err := cli.GetContainerInfo(context.Background(), id)
	if err != nil {
		return nil, err
	}

	imageInfo, err := cli.GetImageInfo(context.Background(), info.Image)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"container": info,
		"image":     imageInfo,
	}, nil
}

func (a runnerDockerAdapter) CheckForUpdates(ctx context.Context, c *types.Container) error {
	imageName := c.GetImageNameWithTag()

	res, err := a.pullImage(ctx, imageName)
	if err != nil {
		return err
	}
	defer res.Close()

	client := containersapi.NewContainersKernelClient(ctx)
	imageInfo, err := client.GetImageInfo(context.Background(), imageName)
	if err != nil {
		return err
	}

	latestImageID := imageInfo.ID

	currentImageID, err := a.getImageID(ctx, *c)
	if err != nil {
		return err
	}

	if latestImageID == currentImageID {
		log.Info("already up-to-date",
			vlog.String("uuid", c.ID.String()),
		)
		c.Update = nil
	} else {
		log.Info("a new update is available",
			vlog.String("uuid", c.ID.String()),
		)
		c.Update = &types.ContainerUpdate{
			CurrentVersion: currentImageID,
			LatestVersion:  latestImageID,
		}
	}

	return nil
}

func (a runnerDockerAdapter) GetAllVersions(ctx context.Context, c types.Container) ([]string, error) {
	log.Info("querying all versions of image", vlog.String("image", c.Image))
	return crane.ListTags(c.Image)
}

func (a runnerDockerAdapter) HasUpdateAvailable(ctx context.Context, c types.Container) (bool, error) {
	//TODO implement me
	return false, nil
}

func (a runnerDockerAdapter) WaitCondition(ctx context.Context, c *types.Container, cond types.WaitContainerCondition) error {
	id, err := a.getContainerID(ctx, *c)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	return cli.WaitContainer(context.Background(), id, string(cond))
}

func (a runnerDockerAdapter) getContainer(ctx context.Context, c types.Container) (types.DockerContainer, error) {
	cli := containersapi.NewContainersKernelClient(ctx)
	containers, err := cli.GetContainers(context.Background())
	if err != nil {
		return types.DockerContainer{}, err
	}

	for _, dc := range containers {
		name := dc.Names[0]
		if name == "/"+c.DockerContainerName() {
			return dc, nil
		}
	}

	return types.DockerContainer{}, errors.NotFoundf("docker container")
}

func (a runnerDockerAdapter) getContainerID(ctx context.Context, c types.Container) (string, error) {
	dc, err := a.getContainer(ctx, c)
	if err != nil {
		return "", err
	}
	return dc.ID, nil
}

func (a runnerDockerAdapter) getImageID(ctx context.Context, c types.Container) (string, error) {
	dc, err := a.getContainer(ctx, c)
	if err != nil {
		return "", err
	}
	return dc.ImageID, nil
}

func (a runnerDockerAdapter) pullImage(ctx context.Context, imageName string) (io.ReadCloser, error) {
	options := types.PullImageOptions{Image: imageName}

	cli := containersapi.NewContainersKernelClient(ctx)
	return cli.PullImage(context.Background(), options)
}

func (a runnerDockerAdapter) buildImageFromName(ctx context.Context, imageName string) (io.ReadCloser, error) {
	res, err := a.pullImage(ctx, imageName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a runnerDockerAdapter) createContainer(ctx context.Context, options types.CreateContainerOptions) (string, error) {
	cli := containersapi.NewContainersKernelClient(ctx)
	res, err := cli.CreateContainer(context.Background(), options)
	if err != nil {
		return "", err
	}

	for _, warn := range res.Warnings {
		log.Warn("warning while creating container", vlog.String("warning", warn))
	}
	return res.ID, nil
}

func (a runnerDockerAdapter) readLogs(ctx context.Context, id string) (stdout io.ReadCloser, stderr io.ReadCloser, _ error) {
	cli := containersapi.NewContainersKernelClient(ctx)

	stdout, err := cli.GetContainerStdout(context.Background(), id)
	if err != nil {
		return nil, nil, err
	}

	stderr, err = cli.GetContainerStderr(context.Background(), id)
	if err != nil {
		stdout.Close()
		return nil, nil, err
	}

	return stdout, stderr, nil
}
