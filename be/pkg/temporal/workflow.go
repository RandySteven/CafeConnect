package temporal_client

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	ActivityExecutionInfo struct {
		ActivityName string
		SignalName   string
	}

	WorkflowExecutionData struct {
		ID                     uint64
		WorkflowID             string
		CurrState              string
		RunID                  string
		SignalEvent            string
		activityExecutionInfos []ActivityExecutionInfo
		StartedAt              time.Time
		CompletedAt            time.Time

		temporalClient Temporal
	}
)

// Execute runs all registered activities sequentially, threading the given state
// through each one. Each activity receives the state as input and returns the
// updated state as output. The state must be a non-nil pointer to a serializable struct.
func (w *WorkflowExecutionData) Execute(ctx workflow.Context, state interface{}) error {
	for _, info := range w.activityExecutionInfos {
		future := workflow.ExecuteActivity(ctx, info.ActivityName, state)
		if err := future.Get(ctx, state); err != nil {
			return fmt.Errorf("activity %s failed: %w", info.ActivityName, err)
		}

		if info.SignalName != "" {
			if err := ExecuteChildWorkflow(ctx, info.SignalName, state); err != nil {
				return fmt.Errorf("child workflow for activity %s failed: %w", info.ActivityName, err)
			}
		}
	}
	return nil
}

// AddTransitionActivity registers an activity with the Temporal worker and adds it
// to the execution pipeline. Activities are executed in the order they are added.
func (w *WorkflowExecutionData) AddTransitionActivity(activityName string, signalName string, activityFn interface{}) {
	w.temporalClient.RegisterActivity(ActivityDefinition{
		Name: activityName,
		Fn:   activityFn,
	})

	w.activityExecutionInfos = append(w.activityExecutionInfos, ActivityExecutionInfo{
		ActivityName: activityName,
		SignalName:   signalName,
	})
}

func (w *WorkflowExecutionData) RegisterWorkflow(name string, fn interface{}) {
	w.temporalClient.RegisterWorkflow(WorkflowDefinition{
		Name: name,
		Fn:   fn,
	})
}

func (w *WorkflowExecutionData) GetWorkflowExecutionData(wfCtx workflow.Context, runID string, result interface{}) error {
	err := w.temporalClient.GetWorkflowResult(context.Background(), w.WorkflowID, runID, result)
	if err != nil {
		return fmt.Errorf("failed to get workflow execution data: %w", err)
	}
	return nil
}

func ExecuteChildWorkflow(ctx workflow.Context, signalName string, request interface{}) error {
	childWorkflowRun := workflow.ExecuteChildWorkflow(ctx, "ChildWorkflow")
	var workflowExecution workflow.Execution
	if err := childWorkflowRun.GetChildWorkflowExecution().Get(ctx, &workflowExecution); err != nil {
		return fmt.Errorf("failed to get child workflow execution: %w", err)
	}

	sigFuture := workflow.SignalExternalWorkflow(ctx, workflowExecution.ID, workflowExecution.RunID, signalName, request)
	if err := sigFuture.Get(ctx, nil); err != nil {
		return fmt.Errorf("failed to signal child workflow: %w", err)
	}

	return nil
}

func NewWorkflowExecutionData(
	temporalClient Temporal,
) WorkflowExecutionData {
	return WorkflowExecutionData{
		activityExecutionInfos: make([]ActivityExecutionInfo, 0),
		temporalClient:         temporalClient,
	}
}
