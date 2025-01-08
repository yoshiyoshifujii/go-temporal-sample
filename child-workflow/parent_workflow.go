package child_workflow

import "go.temporal.io/sdk/workflow"

func SampleParentWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: "ABC-SIMPLE_CHILD_WORKFLOW-ID",
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, SampleChildWorkflow, "World")
	var result string
	err := childWorkflowFuture.Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)
		return "", err
	}

	logger.Info("Parent execution completed.", "Result", result)
	return result, nil
}
