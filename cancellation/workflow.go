package cancellation

import (
	"errors"
	"fmt"
	"go.temporal.io/sdk/workflow"
	"time"
)

func YourWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
		HeartbeatTimeout:    5 * time.Second,
		WaitForCancellation: true,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("cancel workflow started")
	var a *Activities // Used to call Activities by function pointer
	defer func() {

		if !errors.Is(ctx.Err(), workflow.ErrCanceled) {
			return
		}

		// When the Workflow is canceled, it has to get a new disconnected context to execute any Activities
		newCtx, _ := workflow.NewDisconnectedContext(ctx)
		future := workflow.ExecuteActivity(newCtx, a.CleanupActivity)
		err := future.Get(newCtx, nil)
		if err != nil {
			logger.Error("CleanupActivity failed.", "Error", err)
		}
	}()

	var result string
	canceledFuture := workflow.ExecuteActivity(ctx, a.ActivityToBeCanceled)
	err := canceledFuture.Get(ctx, &result)
	logger.Info(fmt.Sprintf("ActivityToBeCanceled returned %v, %v", result, err))

	skippedFuture := workflow.ExecuteActivity(ctx, a.ActivityToBeSkipped)
	err = skippedFuture.Get(ctx, nil)
	logger.Error("Error from ActivityToBeSkipped", err)

	logger.Info("Workflow Execution complete.")

	return nil
}
