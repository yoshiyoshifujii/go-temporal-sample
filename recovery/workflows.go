package recovery

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func RecoverWorkflow(ctx workflow.Context, params Params) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Recover workflow started.")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		HeartbeatTimeout:    30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result ListOpenExecutionsResult
	future := workflow.ExecuteActivity(ctx, ListOpenExecutions, params.Type)
	err := future.Get(ctx, &result)
	if err != nil {
		logger.Error("Failed to list open workflow executions.", "Error", err)
		return err
	}

	concurrency := 1
	if params.Concurrency > 0 {
		concurrency = params.Concurrency
	}

	if result.Count < concurrency {
		concurrency = result.Count
	}

	batchSize := result.Count / concurrency
	if result.Count%concurrency != 0 {
		batchSize++
	}

	info := workflow.GetInfo(ctx)
	expiration := info.WorkflowExecutionTimeout
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2,
		MaximumInterval:    10 * time.Second,
		MaximumAttempts:    100,
	}
	ao = workflow.ActivityOptions{
		StartToCloseTimeout: expiration,
		HeartbeatTimeout:    30 * time.Second,
		RetryPolicy:         retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	doneCh := workflow.NewChannel(ctx)
	for i := 0; i < concurrency; i++ {
		startIndex := i * batchSize

		workflow.Go(ctx, func(ctx workflow.Context) {
			future := workflow.ExecuteActivity(ctx, RecoverExecutions, result.ID, startIndex, batchSize)
			err = future.Get(ctx, nil)
			if err != nil {
				logger.Error("Recover executions failed.", "StartIndex", startIndex, "Error", err)
			} else {
				logger.Info("Recover executions completed.", "StartIndex", startIndex)
			}
			doneCh.Send(ctx, "done")
		})
	}

	for i := 0; i < concurrency; i++ {
		doneCh.Receive(ctx, nil)
	}

	logger.Info("Workflow completed.", "Result", result.Count)

	return nil
}

func TripWorkflow(ctx workflow.Context, state UserState) error {
	logger := workflow.GetLogger(ctx)
	workflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	logger.Info(
		"Trip workflow started for User.",
		"User", workflowID,
		"TripeCounter", state.TripCounter,
	)

	err := workflow.SetQueryHandler(ctx, QueryName, func(input []byte) (int, error) {
		return state.TripCounter, nil
	})
	if err != nil {
		logger.Error("SetQueryHandler failed.", "Error", err)
		return err
	}

	tripCh := workflow.GetSignalChannel(ctx, TripSignalName)
	for i := 0; i < 10; i++ {
		var trip TripEvent
		tripCh.Receive(ctx, &trip)
		logger.Info("Trip complete event received.", "ID", trip.ID, "Total", trip.Total)
		state.TripCounter++
	}

	logger.Info("Starting a new run.", "TripCounter", state.TripCounter)
	return workflow.NewContinueAsNewError(ctx, "TripWorkflow", state)
}
