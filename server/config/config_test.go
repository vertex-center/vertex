package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) TestNew() {
	cfg := New()

	suite.Equal(ProductionMode, cfg.mode)
}

func (suite *ConfigTestSuite) TestNewDebug() {
	suite.T().Setenv("DEBUG", "1")
	cfg := New()

	suite.Equal(DebugMode, cfg.mode)
}
