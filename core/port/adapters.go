package port

import (
	"context"
	types2 "github.com/docker/docker/api/types"
	"github.com/vertex-center/vertex/core/types"
	"io"
)

type (
	BaselinesAdapter interface {
		// GetLatest returns the latest available Baseline. This
		// will typically fetch the latest Baseline from a remote source.
		GetLatest(ctx context.Context, channel types.SettingsUpdatesChannel) (types.Baseline, error)
	}

	DockerAdapter interface {
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
		BuildImage(options types.BuildImageOptions) (types2.ImageBuildResponse, error)
	}

	SettingsAdapter interface {
		GetSettings() types.Settings
		GetNotificationsWebhook() *string
		SetNotificationsWebhook(webhook string) error
		GetChannel() *types.SettingsUpdatesChannel
		SetChannel(channel types.SettingsUpdatesChannel) error
	}

	SshAdapter interface {
		GetAll() ([]types.PublicKey, error)
		Add(key string) error
		Remove(fingerprint string) error
	}
)
