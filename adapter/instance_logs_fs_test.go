package adapter

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/types"
)

type InstanceLoggerTestSuite struct {
	suite.Suite

	logger *InstanceLogger
}

func TestInstanceLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(InstanceLoggerTestSuite))
}

func (suite *InstanceLoggerTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "*_logs_test")
	suite.NoError(err)

	suite.logger = &InstanceLogger{
		uuid:   uuid.New(),
		buffer: []types.LogLine{},
		dir:    dir,
	}
}

func (suite *InstanceLoggerTestSuite) TearDownTest() {
	err := os.RemoveAll(suite.logger.dir)
	suite.NoError(err)
}

func (suite *InstanceLoggerTestSuite) TestOpenClose() {
	// Open
	err := suite.logger.Open()
	suite.NoError(err)

	filename := suite.logger.file.Name()
	suite.FileExists(filename)

	// Close
	err = suite.logger.Close()
	suite.NoError(err)

	// Check that the file still exists
	suite.FileExists(filename)
	suite.Nil(suite.logger.file)
}

func (suite *InstanceLoggerTestSuite) TestCron() {
	// Start cron
	err := suite.logger.startCron()
	suite.NoError(err)

	// Check that the cron is running
	suite.NotNil(suite.logger.scheduler)
	suite.True(suite.logger.scheduler.IsRunning())
	suite.Equal(1, suite.logger.scheduler.Len())

	// Stop cron
	err = suite.logger.stopCron()
	suite.NoError(err)

	// Check that the cron is not running
	suite.False(suite.logger.scheduler.IsRunning())
}

type InstanceLogsFSAdapterTestSuite struct {
	suite.Suite

	adapter *InstanceLogsFSAdapter
}

func TestInstanceLogsFSAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(InstanceLogsFSAdapterTestSuite))
}

func (suite *InstanceLogsFSAdapterTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "*_logs_test")
	suite.NoError(err)

	suite.adapter = NewInstanceLogsFSAdapter(&InstanceLogsFSAdapterParams{
		InstancesPath: dir,
	}).(*InstanceLogsFSAdapter)
}

func (suite *InstanceLogsFSAdapterTestSuite) TearDownTest() {
	err := suite.adapter.UnregisterAll()
	suite.NoError(err)

	err = os.RemoveAll(suite.adapter.instancesPath)
	suite.NoError(err)
}

func (suite *InstanceLogsFSAdapterTestSuite) TestRegisterUnregister() {
	instID := uuid.New()

	// Register
	err := suite.adapter.Register(instID)
	suite.NoError(err)

	defer func() {
		// Unregister
		err = suite.adapter.Unregister(instID)
		suite.NoError(err)

		// Check that the logger was unregistered
		suite.Len(suite.adapter.loggers, 0)
		suite.NotContains(suite.adapter.loggers, instID)
	}()

	// Check that the logger was registered
	suite.Len(suite.adapter.loggers, 1)
	suite.Contains(suite.adapter.loggers, instID)

	l, err := suite.adapter.getLogger(instID)
	suite.NoError(err)
	suite.NotNil(l)
	suite.Equal(instID, l.uuid)
	suite.Equal(l.scheduler.Len(), 1)
}

func (suite *InstanceLogsFSAdapterTestSuite) TestUnregisterAll() {
	instID1 := uuid.New()
	instID2 := uuid.New()

	// Register
	err := suite.adapter.Register(instID1)
	suite.NoError(err)

	err = suite.adapter.Register(instID2)
	suite.NoError(err)

	// Unregister
	err = suite.adapter.UnregisterAll()
	suite.NoError(err)

	// Check that the loggers were unregistered
	suite.Len(suite.adapter.loggers, 0)
	suite.NotContains(suite.adapter.loggers, instID1)
	suite.NotContains(suite.adapter.loggers, instID2)
}

func (suite *InstanceLogsFSAdapterTestSuite) TestPush() {
	instID := uuid.New()

	// Register
	err := suite.adapter.Register(instID)
	suite.NoError(err)
	defer func() {
		err := suite.adapter.Unregister(instID)
		suite.NoError(err)
	}()

	// Push
	suite.adapter.Push(instID, types.LogLine{
		Kind: types.LogKindVertexErr,
		Message: &types.LogLineMessageString{
			Value: "test",
		},
	})

	// Check that the log was pushed
	l, err := suite.adapter.getLogger(instID)
	suite.NoError(err)
	suite.Len(l.buffer, 1)
	suite.Equal(types.LogKindVertexErr, l.buffer[0].Kind)
	suite.Equal("test", l.buffer[0].Message.(*types.LogLineMessageString).Value)
}

func (suite *InstanceLogsFSAdapterTestSuite) TestPop() {
	instID := uuid.New()

	// Register
	err := suite.adapter.Register(instID)
	suite.NoError(err)
	defer func() {
		err := suite.adapter.Unregister(instID)
		suite.NoError(err)
	}()

	// Push
	suite.adapter.Push(instID, types.LogLine{
		Kind: types.LogKindVertexErr,
		Message: &types.LogLineMessageString{
			Value: "test",
		},
	})

	// Pop
	line, err := suite.adapter.Pop(instID)
	suite.NoError(err)
	suite.Equal(types.LogKindVertexErr, line.Kind)
	suite.Equal("test", line.Message.(*types.LogLineMessageString).Value)

	// Check that the log was popped
	l, err := suite.adapter.getLogger(instID)
	suite.NoError(err)
	suite.Len(l.buffer, 0)
}
