package retryactivity

import (
	"context"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"time"
)

func BatchProcessingActivity(ctx context.Context, firstTaskID, batchSize int, processDelay time.Duration) error {
	logger := activity.GetLogger(ctx)

	currentTaskID := firstTaskID
	if activity.HasHeartbeatDetails(ctx) {
		// We are retry from a failed attempt, and there is reported progress that we should resume from.
		var completedIdx int
		if err := activity.GetHeartbeatDetails(ctx, &completedIdx); err == nil {
			currentTaskID = completedIdx + 1
			logger.Info("Resuming from failed attempt", "ReportedProgress", completedIdx)
		}
	}

	taskProcessedInThisAttempt := 0 // used to determine when to fail (simulate failure)
	for ; currentTaskID < firstTaskID+batchSize; currentTaskID++ {
		// process task currentTaskID
		logger.Info("processing task", "TaskID", currentTaskID)
		time.Sleep(processDelay) // simulate time spend on processing each task
		activity.RecordHeartbeat(ctx, currentTaskID)
		taskProcessedInThisAttempt++

		// simulate failure after process 1/3 of the tasks
		if taskProcessedInThisAttempt >= batchSize/3 && currentTaskID < firstTaskID+batchSize-1 {
			logger.Info("Activity failed, will retry...")
			// Activity could return *ApplicationError which is always retryable.
			// To return non-retryable error use temporal.NewNonRetryableApplicationError() constructor.
			return temporal.NewApplicationError("some retryable error", "SomeType")
		}
	}

	logger.Info("Activity succeed.")
	return nil
}
