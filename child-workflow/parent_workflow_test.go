package child_workflow

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
	"testing"
)

func Test_Workflow_Integration(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(SampleChildWorkflow)

	env.ExecuteWorkflow(SampleParentWorkflow)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Hello World!", result)
}

func Test_Workflow_Isolation(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(SampleChildWorkflow)
	env.OnWorkflow("SampleChildWorkflow", mock.Anything, mock.Anything).Return(
		func(ctx workflow.Context, name string) (string, error) {
			return "Hi " + name + "!", nil
		},
	)

	env.ExecuteWorkflow(SampleParentWorkflow)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Hi World!", result)
}
