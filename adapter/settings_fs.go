package adapter

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type SettingsFSAdapter struct {
	settingsPath string
	settings     types.Settings
}

type SettingsFSAdapterParams struct {
	settingsPath string
}

func NewSettingsFSAdapter(params *SettingsFSAdapterParams) types.SettingsAdapterPort {
	if params == nil {
		params = &SettingsFSAdapterParams{}
	}
	if params.settingsPath == "" {
		params.settingsPath = path.Join(storage.Path, "settings")
	}

	err := os.MkdirAll(params.settingsPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Default.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", params.settingsPath),
		)
		os.Exit(1)
	}

	adapter := &SettingsFSAdapter{
		settingsPath: params.settingsPath,
	}

	err = adapter.read()
	if err != nil {
		log.Default.Error(err)
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

func (a *SettingsFSAdapter) read() error {
	p := path.Join(a.settingsPath, "settings.json")
	file, err := os.ReadFile(p)

	if errors.Is(err, os.ErrNotExist) {
		log.Default.Info("settings.json doesn't exists or could not be found")
		return nil
	} else if err != nil {
		return err
	}

	err = json.Unmarshal(file, &a.settings)
	if err != nil {
		return err
	}

	return nil
}

func (a *SettingsFSAdapter) write() error {
	p := path.Join(a.settingsPath, "settings.json")

	bytes, err := json.MarshalIndent(a.settings, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(p, bytes, os.ModePerm)
}
