package recovery

import "go.temporal.io/sdk/client"

type (
	Params struct {
		ID          string
		Type        string
		Concurrency int
	}

	ListOpenExecutionsResult struct {
		ID     string
		Count  int
		HostID string
	}

	RestartParams struct {
		Options client.StartWorkflowOptions
		State   UserState
	}

	SignalParams struct {
		Name string
		Data TripEvent
	}

	UserState struct {
		TripCounter int
	}

	TripEvent struct {
		ID    string
		Total int
	}

	ClientKey int
)

const (
	TemporalClientKey ClientKey = iota
	WorkflowExecutionCacheKey
)

const (
	TripSignalName = "trip_event"
	QueryName      = "counter"
)
