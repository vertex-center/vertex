package updater

import (
	"context"
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

func Execute(ctx context.Context, channel baseline.Channel, updaters ...Updater) error {
	latest, err := baseline.Fetch(ctx, channel)
	if err != nil {
		return fmt.Errorf("fetch baseline: %w", err)
	}

	for _, u := range updaters {
		if u.IsInstalled() {
			continue
		}

		currentVersion, err := u.CurrentVersion()
		if err != nil {
			return err
		}

		version, err := latest.GetVersionByID(u.ID())
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
