package saga

import (
	"time"

	"go.uber.org/multierr"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	WorkflowInput struct {
		Activity01Input             ActivityInput
		Activity01CompensationInput ActivityInput
		Activity02Input             ActivityInput
		Activity02CompensationInput ActivityInput
		Activity03Input             ActivityInput
	}
)

func NewWorkflowInput(
	activity01Error bool,
	activity01CompensationError bool,
	activity02Error bool,
	activity02CompensationError bool,
	activity03Error bool,
) WorkflowInput {
	return WorkflowInput{
		Activity01Input: ActivityInput{
			IsError: activity01Error,
		},
		Activity01CompensationInput: ActivityInput{
			IsCompensationError: activity01CompensationError,
		},
		Activity02Input: ActivityInput{
			IsError: activity02Error,
		},
		Activity02CompensationInput: ActivityInput{
			IsCompensationError: activity02CompensationError,
		},
		Activity03Input: ActivityInput{
			IsError: activity03Error,
		},
	}
}

func Workflow(ctx workflow.Context, input WorkflowInput) error {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    3,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	future := workflow.ExecuteActivity(ctx, Activity01, input.Activity01Input)
	err := future.Get(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			future := workflow.ExecuteActivity(ctx, Activity01Compensation, input.Activity01CompensationInput)
			errCompensation := future.Get(ctx, nil)
			err = multierr.Append(err, errCompensation)
		}
	}()

	future = workflow.ExecuteActivity(ctx, Activity02, input.Activity02Input)
	err = future.Get(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			future := workflow.ExecuteActivity(ctx, Activity02Compensation, input.Activity02CompensationInput)
			errCompensation := future.Get(ctx, nil)
			err = multierr.Append(err, errCompensation)
		}
	}()

	future = workflow.ExecuteActivity(ctx, Activity03, input.Activity03Input)
	err = future.Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
