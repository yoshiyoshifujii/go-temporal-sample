package retryactivity

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func RetryWorkflow(ctx workflow.Context) error {
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

	future := workflow.ExecuteActivity(ctx, BatchProcessingActivity, 0, 20, time.Second)
	err := future.Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("workflow completed with error.", "Error", err)
		return err
	}
	workflow.GetLogger(ctx).Info("workflow completed.")
	return nil
}
