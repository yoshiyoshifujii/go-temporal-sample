package retryactivity

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"
	"testing"
	"time"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_Workflow() {
	env := s.NewTestWorkflowEnvironment()
	env.RegisterActivity(BatchProcessingActivity)

	var processedTaskIDs []int
	env.OnActivity(BatchProcessingActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(func(ctx context.Context, firstTaskID, batchSize int, processDelay time.Duration) error {
			currentTaskID := firstTaskID
			if activity.HasHeartbeatDetails(ctx) {
				var completedIdx int
				if err := activity.GetHeartbeatDetails(ctx, &completedIdx); err == nil {
					currentTaskID = completedIdx + 1
				}
			}
			processedTaskIDs = append(processedTaskIDs, currentTaskID)

			return BatchProcessingActivity(ctx, firstTaskID, batchSize, time.Nanosecond)
		})

	env.ExecuteWorkflow(RetryWorkflow)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	s.Equal([]int{0, 6, 12, 18}, processedTaskIDs)
	env.AssertExpectations(s.T())
}
