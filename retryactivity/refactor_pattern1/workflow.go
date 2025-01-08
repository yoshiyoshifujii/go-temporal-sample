package refactor_pattern1

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

type BatchProcessingWorkflow struct {
	BatchProcessor *BatchProcessor
}

// NewBatchProcessingWorkflow コンストラクタで依存を注入
func NewBatchProcessingWorkflow(batchProcessor *BatchProcessor) *BatchProcessingWorkflow {
	return &BatchProcessingWorkflow{BatchProcessor: batchProcessor}
}

// ExecuteWorkflow バッチ処理のWorkflow
func (bpw *BatchProcessingWorkflow) ExecuteWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
		HeartbeatTimeout:    10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	future := workflow.ExecuteActivity(ctx, bpw.BatchProcessor.BatchProcessingActivity)
	err := future.Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("workflow completed with error.", "Error", err)
		return err
	}
	workflow.GetLogger(ctx).Info("workflow completed.")
	return nil
}
