package services

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
)

// TODO: Move webhooks use to a Discord adapter

type NotificationsService struct {
	settingsAdapter types.SettingsAdapterPort
	eventsAdapter   types.EventAdapterPort
	instanceAdapter types.InstanceAdapterPort

	client   webhook.Client
	listener types.Listener
}

func NewNotificationsService(settingsAdapter types.SettingsAdapterPort, eventsAdapter types.EventAdapterPort, instanceAdapter types.InstanceAdapterPort) NotificationsService {
	return NotificationsService{
		settingsAdapter: settingsAdapter,
		eventsAdapter:   eventsAdapter,
		instanceAdapter: instanceAdapter,
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
		case types.EventInstanceStatusChange:
			if e.Status == types.InstanceStatusOff || e.Status == types.InstanceStatusError || e.Status == types.InstanceStatusRunning {
				s.sendStatus(e.InstanceUUID, e.Status)
			}
		}
	})

	s.eventsAdapter.AddListener(s.listener)

	return nil
}

func (s *NotificationsService) StopWebhook() {
	s.eventsAdapter.RemoveListener(s.listener)
}

func (s *NotificationsService) sendStatus(instanceUUID uuid.UUID, status string) {
	var color int

	switch status {
	case types.InstanceStatusRunning:
		color = 5763719
	case types.InstanceStatusOff:
		color = 15548997
	case types.InstanceStatusError:
		color = 10038562
	}

	instance, err := s.instanceAdapter.Get(instanceUUID)
	if err != nil {
		log.Default.Error(err)
		return
	}

	var name = instance.Name
	if instance.DisplayName != nil {
		name = *instance.DisplayName
	}

	embed := discord.NewEmbedBuilder().
		SetTitle(name).
		SetDescriptionf("Status: %s", status).
		SetColor(color).
		Build()

	_, err = s.client.CreateEmbeds([]discord.Embed{embed})
	if err != nil {
		return
	}
}
