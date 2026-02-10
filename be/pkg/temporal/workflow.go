package temporal_client

import (
	"time"
)

type (
	WorkflowExecution struct {
		ID                     uint64
		WorkflowID             string
		RunID                  string
		ActivityName           string
		TransitionActivityName string
		PreviousWorkflowID     string
		Metadata               map[string]interface{}

		StartedAt   time.Time
		CompletedAt time.Time
	}
)
