package child_workflow_continue_as_new

import (
	"errors"
	"fmt"
	"go.temporal.io/sdk/workflow"
)

func SampleChildWorkflow(ctx workflow.Context, totalCount, runCount int) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Child workflow execution started.")
	if runCount <= 0 {
		logger.Error("Invalid run count.", "RunCount", runCount)
		return "", errors.New("invalid run count")
	}

	totalCount++
	runCount--
	if runCount == 0 {
		result := fmt.Sprintf("Child workflow completed after %v runs.", totalCount)
		logger.Info("Child workflow completed.", "Result", result)
		return result, nil
	}

	logger.Info("Child workflow starting new run.", "RunCount", runCount, "TotalCount", totalCount)
	return "", workflow.NewContinueAsNewError(ctx, SampleChildWorkflow, totalCount, runCount)
}
