package refactor_pattern1

import (
	"context"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

type BatchProcessor struct {
	doSomething            DoSomething
	firstTaskID, batchSize int
}

// NewBatchProcessor コンストラクタで依存を注入
func NewBatchProcessor(doSomething DoSomething, firstTaskID, batchSize int) *BatchProcessor {
	return &BatchProcessor{
		doSomething: doSomething,
		firstTaskID: firstTaskID,
		batchSize:   batchSize,
	}
}

// BatchProcessingActivity バッチ処理のActivity
func (bp *BatchProcessor) BatchProcessingActivity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)

	currentTaskID := bp.firstTaskID
	if activity.HasHeartbeatDetails(ctx) {
		var completedIdx int
		if err := activity.GetHeartbeatDetails(ctx, &completedIdx); err == nil {
			currentTaskID = completedIdx + 1
			logger.Info("Resuming from failed attempt", "ReportedProgress", completedIdx)
		}
	}

	for ; currentTaskID < bp.EndTaskID(); currentTaskID++ {
		logger.Info("processing task", "TaskID", currentTaskID)

		// 依存を注入したDoSomethingを呼び出す
		err := bp.doSomething.Do(ctx, currentTaskID)
		activity.RecordHeartbeat(ctx, currentTaskID)
		if err != nil {
			logger.Info("Activity failed, will retry...")
			// Activity could return *ApplicationError which is always retryable.
			// To return non-retryable error use temporal.NewNonRetryableApplicationError() constructor.
			return temporal.NewApplicationError("some retryable error", err.Error(), err)
		}
	}

	logger.Info("Activity succeed.")
	return nil
}

func (bp *BatchProcessor) EndTaskID() int {
	return bp.firstTaskID + bp.batchSize
}
