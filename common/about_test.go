package common

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/common/baseline"
)

type AboutTestSuite struct {
	suite.Suite
}

func TestAboutTestSuite(t *testing.T) {
	suite.Run(t, new(AboutTestSuite))
}

func (suite *AboutTestSuite) TestNewAbout() {
	about := NewAbout("1.2.3", "cd4ba2876f45775287f426c13adb1868f7c96222", "2006-01-02T15:04:05Z07:00")
	suite.Equal("1.2.3", about.Version)
	suite.Equal("cd4ba2876f45775287f426c13adb1868f7c96222", about.Commit)
	suite.Equal("2006-01-02T15:04:05Z07:00", about.Date)
	suite.Equal(runtime.GOOS, about.OS)
	suite.Equal(runtime.GOARCH, about.Arch)
}

func (suite *AboutTestSuite) TestChannel() {
	suite.Equal(baseline.ChannelStable, NewAbout("1.2.3", "", "").Channel())
	suite.Equal(baseline.ChannelBeta, NewAbout("1.2.3-beta", "", "").Channel())
}
