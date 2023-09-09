package services

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
)

// TODO: Move webhooks use to a Discord repo

type NotificationsService struct {
	settingsRepo types.SettingsRepository
	eventsRepo   types.EventRepository
	instanceRepo types.InstanceRepository

	client   webhook.Client
	listener types.Listener
}

func NewNotificationsService(settingsRepo types.SettingsRepository, eventsRepo types.EventRepository, instanceRepo types.InstanceRepository) NotificationsService {
	return NotificationsService{
		settingsRepo: settingsRepo,
		eventsRepo:   eventsRepo,
		instanceRepo: instanceRepo,
	}
}

func (s *NotificationsService) StartWebhook() error {
	webhookURL := s.settingsRepo.GetNotificationsWebhook()
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

	s.eventsRepo.AddListener(s.listener)

	return nil
}

func (s *NotificationsService) StopWebhook() {
	s.eventsRepo.RemoveListener(s.listener)
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

	instance, err := s.instanceRepo.Get(instanceUUID)
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
