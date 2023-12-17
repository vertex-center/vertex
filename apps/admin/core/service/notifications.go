package service

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	apptypes "github.com/vertex-center/vertex/common/app"
	ev "github.com/vertex-center/vertex/common/event"
	"github.com/vertex-center/vertex/common/uuid"
	"github.com/vertex-center/vertex/pkg/event"
)

// TODO: Move webhooks use to a Discord adapter

type notificationsService struct {
	uuid            uuid.UUID
	ctx             *apptypes.Context
	settingsAdapter port.SettingsAdapter
	client          webhook.Client
}

func NewNotificationsService(ctx *apptypes.Context, settingsAdapter port.SettingsAdapter) notificationsService {
	return notificationsService{
		uuid:            uuid.New(),
		ctx:             ctx,
		settingsAdapter: settingsAdapter,
	}
}

func (s *notificationsService) StartWebhook() error {
	adminSetting, err := s.settingsAdapter.Get()
	if err != nil {
		return err
	}

	url := adminSetting.Webhook
	if url == nil {
		return nil
	}

	s.client, err = webhook.NewWithURL(*url)
	if err != nil {
		return err
	}

	s.ctx.AddListener(s)

	return nil
}

func (s *notificationsService) StopWebhook() {
	s.ctx.RemoveListener(s)
}

func (s *notificationsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *notificationsService) OnEvent(e event.Event) error {
	switch e := e.(type) {
	case ev.ServerSetupCompleted:
		return s.StartWebhook()
	case containerstypes.EventContainerStatusChange:
		if e.Status == containerstypes.ContainerStatusOff || e.Status == containerstypes.ContainerStatusError || e.Status == containerstypes.ContainerStatusRunning {
			s.sendStatus(e.Name, e.Status)
		}
	}
	return nil
}

func (s *notificationsService) sendStatus(name string, status string) {
	var color int

	switch status {
	case containerstypes.ContainerStatusRunning:
		color = 5763719
	case containerstypes.ContainerStatusOff:
		color = 15548997
	case containerstypes.ContainerStatusError:
		color = 10038562
	}

	embed := discord.NewEmbedBuilder().
		SetTitle(name).
		SetDescriptionf("Status: %s", status).
		SetColor(color).
		Build()

	_, err := s.client.CreateEmbeds([]discord.Embed{embed})
	if err != nil {
		return
	}
}
