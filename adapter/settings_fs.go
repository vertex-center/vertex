package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

var (
	errSettingsNotFound       = errors.New("settings.json doesn't exists or could not be found")
	errSettingsFailedToRead   = errors.New("failed to read settings.json")
	errSettingsFailedToDecode = errors.New("failed to decode settings.json")
)

type SettingsFSAdapter struct {
	settingsDir string
	settings    types.Settings
}

type SettingsFSAdapterParams struct {
	settingsDir string
}

func NewSettingsFSAdapter(params *SettingsFSAdapterParams) types.SettingsAdapterPort {
	if params == nil {
		params = &SettingsFSAdapterParams{}
	}
	if params.settingsDir == "" {
		params.settingsDir = path.Join(storage.Path, "settings")
	}

	err := os.MkdirAll(params.settingsDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", params.settingsDir),
		)
		os.Exit(1)
	}

	adapter := &SettingsFSAdapter{
		settingsDir: params.settingsDir,
	}

	err = adapter.read()
	if errors.Is(err, errSettingsFailedToDecode) {
		log.Error(err)
	}

	return adapter
}

func (a *SettingsFSAdapter) GetSettings() types.Settings {
	return a.settings
}

func (a *SettingsFSAdapter) GetNotificationsWebhook() *string {
	if a.settings.Notifications == nil {
		return nil
	}
	return a.settings.Notifications.Webhook
}

func (a *SettingsFSAdapter) SetNotificationsWebhook(webhook *string) error {
	if a.settings.Notifications == nil {
		a.settings.Notifications = &types.SettingsNotifications{}
	}
	a.settings.Notifications.Webhook = webhook
	return a.write()
}

func (a *SettingsFSAdapter) GetChannel() *types.SettingsUpdatesChannel {
	if a.settings.Updates == nil {
		return nil
	}
	return a.settings.Updates.Channel
}

func (a *SettingsFSAdapter) SetChannel(channel *types.SettingsUpdatesChannel) error {
	if a.settings.Updates == nil {
		a.settings.Updates = &types.SettingsUpdates{}
	}
	a.settings.Updates.Channel = channel
	return a.write()
}

func (a *SettingsFSAdapter) read() error {
	p := path.Join(a.settingsDir, "settings.json")
	file, err := os.ReadFile(p)

	if errors.Is(err, fs.ErrNotExist) {
		return errSettingsNotFound
	} else if err != nil {
		return fmt.Errorf("%w: %w", errSettingsFailedToRead, err)
	}

	err = json.Unmarshal(file, &a.settings)
	if err != nil {
		return fmt.Errorf("%w: %w", errSettingsFailedToDecode, err)
	}
	return nil
}

func (a *SettingsFSAdapter) write() error {
	p := path.Join(a.settingsDir, "settings.json")

	bytes, err := json.MarshalIndent(a.settings, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(p, bytes, os.ModePerm)
}
