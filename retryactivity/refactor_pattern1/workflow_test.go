package refactor_pattern1

import (
	"context"
	"fmt"
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

type FailDoSomething struct {
	firstTaskID, batchSize int
}

func (ds *FailDoSomething) Do(ctx context.Context, currentTaskID int) error {
	time.Sleep(time.Nanosecond)
	// simulate failure after process 1/3 of the tasks
	if (currentTaskID-ds.firstTaskID+1)%(ds.batchSize/3) == 0 && currentTaskID < ds.firstTaskID+ds.batchSize-1 {
		return fmt.Errorf("task %d failed", currentTaskID)
	}

	return nil
}

func (s *UnitTestSuite) Test_Workflow() {
	nbp := NewBatchProcessor(&FailDoSomething{
		firstTaskID: 0,
		batchSize:   20,
	}, 0, 20)

	env := s.NewTestWorkflowEnvironment()
	env.RegisterActivity(nbp.BatchProcessingActivity)

	var processedTaskIDs []int
	env.OnActivity(nbp.BatchProcessingActivity, mock.Anything).
		Return(func(ctx context.Context) error {
			currentTaskID := nbp.firstTaskID
			if activity.HasHeartbeatDetails(ctx) {
				var completedIdx int
				if err := activity.GetHeartbeatDetails(ctx, &completedIdx); err == nil {
					currentTaskID = completedIdx + 1
				}
			}
			processedTaskIDs = append(processedTaskIDs, currentTaskID)

			return nbp.BatchProcessingActivity(ctx)
		})

	bpw := NewBatchProcessingWorkflow(nbp)
	env.ExecuteWorkflow(bpw.ExecuteWorkflow)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	s.Equal([]int{0, 6, 12, 18}, processedTaskIDs)
	env.AssertExpectations(s.T())
}
