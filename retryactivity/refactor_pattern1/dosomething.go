package refactor_pattern1

import (
	"context"
	"go.temporal.io/sdk/activity"
	"time"
)

type DoSomething interface {
	Do(ctx context.Context, taskID int) error
}

type DoSomethingImpl struct {
	processDelay time.Duration
}

// NewDoSomethingImpl コンストラクタ
func NewDoSomethingImpl(processDelay time.Duration) *DoSomethingImpl {
	return &DoSomethingImpl{
		processDelay: processDelay,
	}
}

func (dsi *DoSomethingImpl) Do(ctx context.Context, taskID int) error {
	logger := activity.GetLogger(ctx)
	logger.Info("do delay processing", "TaskID", taskID)
	time.Sleep(dsi.processDelay) // simulate time spend on processing each task
	return nil
}
