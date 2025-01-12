package child_workflow_continue_as_new

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
)

func SampleParentWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	info := workflow.GetInfo(ctx)
	execution := info.WorkflowExecution

	// 親ワークフローは子実行のためのIDを指定することができます。
	// 各実行のためにIDが一意であることを確認してください。
	// 子実行のためにTemporalサーバーが一意のIDを生成することを望む場合は指定しないでください。
	childId := fmt.Sprintf("child_workflow:%v", execution.RunID)
	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: childId,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)
	var result string
	future := workflow.ExecuteChildWorkflow(ctx, SampleChildWorkflow, 0, 5)
	err := future.Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)
		return err
	}

	logger.Info("Parent execution completed.", "Result", result)
	return nil
}
