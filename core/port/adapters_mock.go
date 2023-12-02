package port

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vertex-center/vertex/core/types"
)

type (
	MockBaselinesAdapter struct {
		mock.Mock
	}

	MockAdminSettingsAdapter struct {
		mock.Mock
	}
)

func (m *MockBaselinesAdapter) GetLatest(ctx context.Context, channel types.UpdatesChannel) (types.Baseline, error) {
	args := m.Called(ctx, channel)
	return args.Get(0).(types.Baseline), args.Error(1)
}

func (m *MockAdminSettingsAdapter) Get() (types.AdminSettings, error) {
	args := m.Called()
	return args.Get(0).(types.AdminSettings), args.Error(1)
}

func (m *MockAdminSettingsAdapter) Update(settings types.AdminSettings) error {
	args := m.Called(settings)
	return args.Error(0)
}
