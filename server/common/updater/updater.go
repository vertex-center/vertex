package updater

import (
	"fmt"

	"github.com/vertex-center/vertex/common/baseline"
)

type Updater interface {
	CurrentVersion() (string, error)
	Install(version string) error
	IsInstalled() bool
	ID() string
}

func CheckUpdates(bl baseline.Baseline, updaters ...Updater) (bool, error) {
	for _, u := range updaters {
		currentVersion, err := u.CurrentVersion()
		if err != nil {
			return false, err
		}

		latestVersion, err := bl.GetVersionByID(u.ID())
		if err != nil {
			return false, fmt.Errorf("'%w' when accessing '%s'", err, u.ID())
		}

		if currentVersion != latestVersion {
			return true, nil
		}
	}

	return false, nil
}

func Install(bl baseline.Baseline, updaters ...Updater) error {
	var err error
	for _, u := range updaters {
		currentVersion := "not-installed"

		if u.IsInstalled() {
			currentVersion, err = u.CurrentVersion()
			if err != nil {
				return err
			}
		}

		version, err := bl.GetVersionByID(u.ID())
		if err != nil {
			return err
		}

		if version != currentVersion {
			err = u.Install(version)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
