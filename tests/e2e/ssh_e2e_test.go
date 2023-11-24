//go:build e2e

package e2e

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SshE2ETestSuite struct {
	suite.Suite
}

func TestE2ESshTestSuite(t *testing.T) {
	suite.Run(t, new(SshE2ETestSuite))
}
