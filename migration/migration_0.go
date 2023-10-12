package migration

import (
	"encoding/json"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

// Migration 0: migrate instance_settings.json to settings.yml.
// This migration is needed to unify the formats.
type migration0 struct{}

func (m *migration0) renameSettingsFile(instancePath string) error {
	const (
		oldInstanceSettingsPath = ".vertex/instance_settings.json"
		newInstanceSettingsPath = ".vertex/settings.yml"
	)

	settingsJson, err := os.ReadFile(path.Join(instancePath, oldInstanceSettingsPath))
	if err != nil && os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	var v interface{}
	err = json.Unmarshal(settingsJson, &v)
	if err != nil {
		return err
	}

	settingsYaml, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(instancePath, newInstanceSettingsPath), settingsYaml, 0644)
	if err != nil {
		return err
	}

	return os.Remove(path.Join(instancePath, oldInstanceSettingsPath))
}

func (m *migration0) Up(livePath string) error {
	instancesPath := path.Join(livePath, "instances")

	dirs, err := os.ReadDir(instancesPath)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		err := m.renameSettingsFile(path.Join(instancesPath, dir.Name()))
		if err != nil {
			continue
		}
	}

	return nil
}
