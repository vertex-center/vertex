package service

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/event"
)

// TODO: Move webhooks use to a Discord adapter

type NotificationsService struct {
	uuid            uuid.UUID
	ctx             *apptypes.Context
	settingsAdapter port.AdminSettingsAdapter
	client          webhook.Client
}

func NewNotificationsService(ctx *apptypes.Context, settingsAdapter port.AdminSettingsAdapter) NotificationsService {
	return NotificationsService{
		uuid:            uuid.New(),
		ctx:             ctx,
		settingsAdapter: settingsAdapter,
	}
}

func (s *NotificationsService) StartWebhook() error {
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

func (s *NotificationsService) StopWebhook() {
	s.ctx.RemoveListener(s)
}

func (s *NotificationsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *NotificationsService) OnEvent(e event.Event) error {
	switch e := e.(type) {
	case types.EventServerSetupCompleted:
		return s.StartWebhook()
	case containerstypes.EventContainerStatusChange:
		if e.Status == containerstypes.ContainerStatusOff || e.Status == containerstypes.ContainerStatusError || e.Status == containerstypes.ContainerStatusRunning {
			s.sendStatus(e.Name, e.Status)
		}
	}
	return nil
}

func (s *NotificationsService) sendStatus(name string, status string) {
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