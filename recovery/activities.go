package recovery

import (
	"context"
	"errors"
	"github.com/pborman/uuid"
	"github.com/yoshiyoshifujii/go-temporal-sample/recovery/cache"
	commonpb "go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"
	filterpb "go.temporal.io/api/filter/v1"
	historypb "go.temporal.io/api/history/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// HostID - Use a new uuid just for demo so we can run 2 host specific activity workers on same machine.
// In real world case, you would use a hostname or ip address as HostID.
var HostID = "recovery_" + uuid.New()

var (
	ErrClientNotFound         = errors.New("failed to retrieve client from context")
	ErrExecutionCacheNotFound = errors.New("failed to retrieve cache from context")
)

func ListOpenExecutions(ctx context.Context, workflowType string) (*ListOpenExecutionsResult, error) {
	key := uuid.New()
	logger := activity.GetLogger(ctx)
	logger.Info(
		"List all open executions of type.",
		"WorkflowType", workflowType,
		"HostID", HostID,
	)

	c, err := getClientFromContext(ctx)
	if err != nil {
		return nil, err
	}

	executionCache := ctx.Value(WorkflowExecutionCacheKey).(cache.Cache)
	if executionCache == nil {
		logger.Error("Could not retrieve cache from context.")
		return nil, ErrExecutionCacheNotFound
	}

	openExecutions, err := getAllExecutionsOfType(ctx, c, workflowType)
	if err != nil {
		return nil, err
	}

	executionCache.Put(key, openExecutions)
	return &ListOpenExecutionsResult{
		ID:     key,
		Count:  len(openExecutions),
		HostID: HostID,
	}, nil
}

func getAllExecutionsOfType(ctx context.Context, c client.Client, workflowType string) ([]*commonpb.WorkflowExecution, error) {
	var openExecutions []*commonpb.WorkflowExecution
	var nextPageToken []byte
	for hasMore := true; hasMore; hasMore = len(nextPageToken) > 0 {
		zeroTime := time.Time{}
		now := time.Now()
		resp, err := c.ListOpenWorkflow(ctx, &workflowservice.ListOpenWorkflowExecutionsRequest{
			Namespace:       client.DefaultNamespace,
			MaximumPageSize: 10,
			NextPageToken:   nextPageToken,
			StartTimeFilter: &filterpb.StartTimeFilter{
				EarliestTime: timestamppb.New(zeroTime),
				LatestTime:   timestamppb.New(now),
			},
			Filters: &workflowservice.ListOpenWorkflowExecutionsRequest_TypeFilter{
				TypeFilter: &filterpb.WorkflowTypeFilter{
					Name: workflowType,
				},
			},
		})
		if err != nil {
			return nil, err
		}

		for _, r := range resp.Executions {
			openExecutions = append(openExecutions, r.Execution)
		}

		nextPageToken = resp.NextPageToken
		activity.RecordHeartbeat(ctx, nextPageToken)
	}
	return openExecutions, nil
}

func getClientFromContext(ctx context.Context) (client.Client, error) {
	logger := activity.GetLogger(ctx)
	temporalClient := ctx.Value(TemporalClientKey).(client.Client)
	if temporalClient == nil {
		logger.Error("Could not retrieve temporal client from context.")
		return nil, ErrClientNotFound
	}
	return temporalClient, nil
}

func RecoverExecutions(ctx context.Context, key string, startIndex, batchSize int) error {
	logger := activity.GetLogger(ctx)
	logger.Info(
		"Starting execution recovery.",
		"HostID", HostID,
		"Key", key,
		"StartIndex", startIndex,
		"BatchSize", batchSize,
	)

	executionCache := ctx.Value(WorkflowExecutionCacheKey).(cache.Cache)
	if executionCache == nil {
		logger.Error("Could not retrieve cache from context.")
		return ErrExecutionCacheNotFound
	}

	openExecutions := executionCache.Get(key).([]*commonpb.WorkflowExecution)
	endIndex := startIndex + batchSize

	if activity.HasHeartbeatDetails(ctx) {
		var finishedIndex int
		if err := activity.GetHeartbeatDetails(ctx, &finishedIndex); err == nil {
			startIndex = finishedIndex + 1
		}
	}

	for index := startIndex; index < endIndex && index < len(openExecutions); index++ {
		execution := openExecutions[index]
		if err := recoverSingleExecution(ctx, execution.GetWorkflowId()); err != nil {
			logger.Error(
				"Failed to recover execution.",
				"WorkflowID", execution.GetWorkflowId(),
				"Error", err,
			)
			return err
		}

		activity.RecordHeartbeat(ctx, index)
	}
	return nil
}

func recoverSingleExecution(ctx context.Context, workflowID string) error {
	logger := activity.GetLogger(ctx)
	c, err := getClientFromContext(ctx)
	if err != nil {
		return err
	}

	execution := &commonpb.WorkflowExecution{
		WorkflowId: workflowID,
	}
	history, err := getHistory(ctx, execution)
	if err != nil {
		return err
	}

	if len(history) == 0 {
		return nil
	}

	firstEvent := history[0]
	lastEvent := history[len(history)-1]

	params, err := extractStateFromEvent(workflowID, firstEvent)
	if err != nil {
		return err
	}

	signals, err := extractSignals(history)
	if err != nil {
		return err
	}

	if !isExecutionCompleted(lastEvent) {
		err := c.TerminateWorkflow(ctx, execution.GetWorkflowId(), execution.GetRunId(), "Recover", nil)
		if err != nil {
			return err
		}
	}

	newRun, err := c.ExecuteWorkflow(ctx, params.Options, "TripWorkflow", params.State)
	if err != nil {
		return err
	}

	for _, s := range signals {
		_ = c.SignalWorkflow(ctx, execution.GetWorkflowId(), newRun.GetRunID(), s.Name, s.Data)
	}

	logger.Info(
		"Successfully restarted workflow.",
		"WorkflowID", execution.GetWorkflowId(),
		"NewRunID", newRun.GetRunID(),
	)

	return nil
}

func isExecutionCompleted(event *historypb.HistoryEvent) bool {
	switch event.GetEventType() {
	case enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_COMPLETED,
		enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_TERMINATED,
		enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_CANCELED,
		enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_FAILED,
		enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_TIMED_OUT:
		return true
	default:
		return false
	}
}

func extractSignals(events []*historypb.HistoryEvent) ([]*SignalParams, error) {
	var signals []*SignalParams
	for _, event := range events {
		if event.GetEventType() == enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_SIGNALED {
			attr := event.GetWorkflowExecutionSignaledEventAttributes()
			if attr.GetSignalName() == TripSignalName && attr.GetInput() != nil {
				signalData, err := deserializeTripEvent(attr.GetInput())
				if err != nil {
					return nil, err
				}

				signal := &SignalParams{
					Name: attr.GetSignalName(),
					Data: signalData,
				}
				signals = append(signals, signal)
			}
		}
	}
	return signals, nil
}

func deserializeTripEvent(data *commonpb.Payloads) (TripEvent, error) {
	var trip TripEvent
	if err := converter.GetDefaultDataConverter().FromPayloads(data, &trip); err != nil {
		return TripEvent{}, err
	}
	return trip, nil
}

func extractStateFromEvent(workflowID string, event *historypb.HistoryEvent) (*RestartParams, error) {
	switch event.GetEventType() {
	case enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_STARTED:
		attr := event.GetWorkflowExecutionStartedEventAttributes()
		state, err := deserializeUserState(attr.GetInput())
		if err != nil {
			return nil, err
		}
		return &RestartParams{
			Options: client.StartWorkflowOptions{
				ID:                  workflowID,
				TaskQueue:           attr.TaskQueue.GetName(),
				WorkflowTaskTimeout: attr.GetWorkflowTaskTimeout().AsDuration(),
			},
			State: state,
		}, nil
	default:
		return nil, errors.New("unknown event type")
	}
}

func deserializeUserState(data *commonpb.Payloads) (UserState, error) {
	var state UserState
	if err := converter.GetDefaultDataConverter().FromPayloads(data, &state); err != nil {
		return UserState{}, err
	}
	return state, nil
}

func getHistory(ctx context.Context, execution *commonpb.WorkflowExecution) ([]*historypb.HistoryEvent, error) {
	return nil, nil
}
