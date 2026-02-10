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

// func (w *WorkflowExecution) GetID() string {
// 	return w.ID
// }

// func (w *WorkflowExecution) GetRunID() string {
// 	return w.RunID
// }

// func (w *WorkflowExecution) GetActivityName() string {
// 	return w.ActivityName
// }

// func (w *WorkflowExecution) GetTransitionActivityName() string {
// 	return w.TransitionActivityName
// }

// func (w *WorkflowExecution) GetPreviousWorkflowID() string {
// 	return w.PreviousWorkflowID
// }
