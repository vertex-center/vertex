package services

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/types"
)

// TODO: Move webhooks use to a Discord adapter

type NotificationsService struct {
	ctx             *types.VertexContext
	settingsAdapter types.SettingsAdapterPort
	client          webhook.Client
	listener        types.Listener
}

func NewNotificationsService(ctx *types.VertexContext, settingsAdapter types.SettingsAdapterPort) NotificationsService {
	return NotificationsService{
		ctx:             ctx,
		settingsAdapter: settingsAdapter,
	}
}

func (s *NotificationsService) StartWebhook() error {
	webhookURL := s.settingsAdapter.GetNotificationsWebhook()
	if webhookURL == nil {
		return nil
	}

	var err error
	s.client, err = webhook.NewWithURL(*webhookURL)
	if err != nil {
		return err
	}

	s.listener = types.NewTempListener(func(e interface{}) {
		switch e := e.(type) {
		case instancestypes.EventInstanceStatusChange:
			if e.Status == instancestypes.InstanceStatusOff || e.Status == instancestypes.InstanceStatusError || e.Status == instancestypes.InstanceStatusRunning {
				s.sendStatus(e.Name, e.Status)
			}
		}
	})

	s.ctx.AddListener(s.listener)

	return nil
}

func (s *NotificationsService) StopWebhook() {
	s.ctx.RemoveListener(s.listener)
}

func (s *NotificationsService) sendStatus(name string, status string) {
	var color int

	switch status {
	case instancestypes.InstanceStatusRunning:
		color = 5763719
	case instancestypes.InstanceStatusOff:
		color = 15548997
	case instancestypes.InstanceStatusError:
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
