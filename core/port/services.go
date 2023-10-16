package port

import (
	dockertypes "github.com/docker/docker/api/types"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"io"
)

type (
	AppsService interface {
		All() []app.Meta
	}

	DockerService interface {
		ListContainers() ([]types.Container, error)
		DeleteContainer(id string) error
		CreateContainer(options types.CreateContainerOptions) (types.CreateContainerResponse, error)
		StartContainer(id string) error
		StopContainer(id string) error
		InfoContainer(id string) (types.InfoContainerResponse, error)
		LogsStdoutContainer(id string) (io.ReadCloser, error)
		LogsStderrContainer(id string) (io.ReadCloser, error)
		WaitContainer(id string, cond types.WaitContainerCondition) error
		InfoImage(id string) (types.InfoImageResponse, error)
		PullImage(options types.PullImageOptions) (io.ReadCloser, error)
		BuildImage(options types.BuildImageOptions) (dockertypes.ImageBuildResponse, error)
	}

	HardwareService interface {
		Get() types.Hardware
	}

	SettingsService interface {
		Get() types.Settings
		Update(settings types.Settings) error
		GetNotificationsWebhook() *string
		SetNotificationsWebhook(webhook string) error
		GetChannel() types.SettingsUpdatesChannel
		SetChannel(channel types.SettingsUpdatesChannel) error
	}

	SshService interface {
		GetAll() ([]types.PublicKey, error)
		Add(key string) error
		Delete(fingerprint string) error
	}

	UpdateService interface {
		GetUpdate(channel types.SettingsUpdatesChannel) (*types.Update, error)
		InstallLatest(channel types.SettingsUpdatesChannel) error
	}
)
