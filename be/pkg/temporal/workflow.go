package temporal_client

import (
	"context"
	"time"
)

type (
	WorkflowExecution interface {
		// GetID returns the workflow ID.
		GetID() string
		// GetRunID returns the run ID.
		GetRunID() string

		// ExecuteLocalActivity executes a local activity.
		ExecuteLocalActivity(ctx context.Context, activityFn interface{}, args ...interface{}) error

		// GetWorkflowResult blocks until the workflow completes and returns the result.
		GetWorkflowResult(ctx context.Context, workflowID string, runID string, result interface{}) error
	}
	WorkflowExecutionInfo struct {
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

// func (t *temporalClient) GetWorkflowRunID(workflowID string, runID string) string {
// 	run := t.client.GetWorkflow(context.Background(), workflowID, runID)
// 	return run.GetRunID()
// }

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
