package saga

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"testing"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}

	// given
	testcases := map[string]struct {
		isErrorCase                 bool
		isActivity01Err             error
		isActivity01CompensationErr error
		isActivity02Err             error
		isActivity02CompensationErr error
		isActivity03Err             error
	}{
		"success": {
			isErrorCase:                 false,
			isActivity01Err:             nil,
			isActivity01CompensationErr: nil,
			isActivity02Err:             nil,
			isActivity02CompensationErr: nil,
			isActivity03Err:             nil,
		},
		"activity01 fails": {
			isErrorCase:                 true,
			isActivity01Err:             assert.AnError,
			isActivity01CompensationErr: nil,
			isActivity02Err:             nil,
			isActivity02CompensationErr: nil,
			isActivity03Err:             nil,
		},
		"activity02 fails, compensation succeeds": {
			isErrorCase:                 true,
			isActivity01Err:             nil,
			isActivity01CompensationErr: nil,
			isActivity02Err:             assert.AnError,
			isActivity02CompensationErr: nil,
			isActivity03Err:             nil,
		},
		"activity02 fails, compensation fails": {
			isErrorCase:                 true,
			isActivity01Err:             nil,
			isActivity01CompensationErr: nil,
			isActivity02Err:             assert.AnError,
			isActivity02CompensationErr: assert.AnError,
			isActivity03Err:             nil,
		},
		"activity03 fails": {
			isErrorCase:                 true,
			isActivity01Err:             nil,
			isActivity01CompensationErr: nil,
			isActivity02Err:             nil,
			isActivity02CompensationErr: nil,
			isActivity03Err:             assert.AnError,
		},
		"activity01 compensation fails": {
			isErrorCase:                 true,
			isActivity01Err:             assert.AnError,
			isActivity01CompensationErr: assert.AnError,
			isActivity02Err:             nil,
			isActivity02CompensationErr: nil,
			isActivity03Err:             nil,
		},
	}

	for name, ts := range testcases {
		t.Run(name, func(t *testing.T) {
			env := testSuite.NewTestWorkflowEnvironment()
			workflowInput := WorkflowInput{}
			activityInput := ActivityInput{}

			// given
			env.OnActivity(Activity01, mock.Anything, activityInput).Return(ts.isActivity01Err)
			env.OnActivity(Activity01Compensation, mock.Anything, activityInput).Return(ts.isActivity01CompensationErr)
			env.OnActivity(Activity02, mock.Anything, activityInput).Return(ts.isActivity02Err)
			env.OnActivity(Activity02Compensation, mock.Anything, activityInput).Return(ts.isActivity02CompensationErr)
			env.OnActivity(Activity03, mock.Anything, activityInput).Return(ts.isActivity03Err)

			// when
			env.ExecuteWorkflow(Workflow, workflowInput)

			// then
			if ts.isErrorCase {
				require.True(t, env.IsWorkflowCompleted())
				require.Error(t, env.GetWorkflowError())
				return
			}

			require.True(t, env.IsWorkflowCompleted())
			require.NoError(t, env.GetWorkflowError())
		})
	}
}
