package expense

import (
	"go.temporal.io/sdk/workflow"
	"time"
)

var (
	expenseServerHostPort = "http://localhost:8099"
)

// SampleExpenseWorkflow workflow definition
func SampleExpenseWorkflow(ctx workflow.Context, expenseID string) (result string, err error) {
	logger := workflow.GetLogger(ctx)

	// step 1, create new expense report
	err = createExpense(ctx, expenseID)
	if err != nil {
		logger.Error("Failed to create expense report", "Error", err)
		return "", err
	}

	// step 2, wait for the expense report to be approved (or rejected)
	status, err := waitForDecision(ctx, expenseID)
	if err != nil {
		return "", err
	}

	if status != "APPROVED" {
		logger.Info("Workflow completed.", "ExpenseStatus", status)
		return "", nil
	}

	// step 3, request payment to the expense
	err = requestPayment(ctx, expenseID)
	if err != nil {
		logger.Error("Workflow completed with payment failed.", "Error", err)
		return "", err
	}

	logger.Info("Workflow completed with expense payment completed.")
	return "COMPLETED", nil
}

func requestPayment(ctx workflow.Context, expenseID string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	future := workflow.ExecuteActivity(ctx, PaymentActivity, expenseID)
	return future.Get(ctx, nil)
}

func waitForDecision(ctx workflow.Context, expenseID string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	// Notice that we set the timeout to be 10 minutes for this sample demo.
	// If the expected time for the activity to complete (waiting for human to approve the request) is longer,
	// you should set the timeout accordingly so the Temporal system will wait accordingly.
	// Otherwise, Temporal system could mark the activity as failure by timeout.
	var status string
	future := workflow.ExecuteActivity(ctx, WaitForDecisionActivity, expenseID)
	err := future.Get(ctx, &status)
	return status, err
}

func createExpense(ctx workflow.Context, expenseID string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	future := workflow.ExecuteActivity(ctx, CreateExpenseActivity, expenseID)
	return future.Get(ctx, nil)
}
