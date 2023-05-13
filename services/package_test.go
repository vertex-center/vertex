package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/types"
)

type PackageTestSuite struct {
	suite.Suite

	service PackageService
	repo    MockPackageRepository
}

func TestEventInMemoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PackageTestSuite))
}

func (suite *PackageTestSuite) SetupSuite() {
	suite.repo = MockPackageRepository{}
	suite.service = NewPackageService(&suite.repo)
}

func (suite *PackageTestSuite) TestInstallationCommand() {
	p := types.Package{
		InstallPackage: map[string]string{
			"brew":    "redis",
			"pacman":  "redis",
			"apt-get": "redis",
			"snap":    "redis",
			"sources": "script:install.sh",
		},
	}

	tests := []struct {
		pm      string
		command string
		sudo    bool
	}{
		{pm: "brew", command: "brew install redis", sudo: false},
		{pm: "pacman", command: "sudo pacman -S --noconfirm redis", sudo: true},
		{pm: "snap", command: "sudo snap install redis", sudo: true},
		{pm: "apt-get", command: "sudo apt-get install redis", sudo: true},
		{pm: "sources", command: "install.sh", sudo: false},
	}

	for _, test := range tests {
		command, err := suite.service.InstallationCommand(&p, test.pm)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), test.command, command.Cmd)
		assert.Equal(suite.T(), test.sudo, command.Sudo)
	}
}

type MockPackageRepository struct {
	mock.Mock
}

func (m *MockPackageRepository) Get(id string) (types.Package, error) {
	m.Called(id)
	return types.Package{}, nil
}

func (m *MockPackageRepository) GetPath(id string) string {
	m.Called(id)
	return ""
}
