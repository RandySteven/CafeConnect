package temporal_client

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type (
	// Navigable allows the Execute state machine to read which activity
	// should run next. Any state struct that implements this interface
	// enables branching in the pipeline.
	//
	// After an activity returns, Execute checks GetNextActivity():
	//   - "" (empty)  → continue to the next sequential activity
	//   - activity name → branch to that activity (must be registered via AddBranchActivity)
	//
	// After branching, execution stops (the branch path runs to completion,
	// then Execute returns). The workflow function can inspect the state to
	// decide what to do next.
	Navigable interface {
		GetNextActivity() string
		SetNextActivity(name string)
	}

	ActivityExecutionInfo struct {
		ActivityName    string
		SignalName      string
		ActivityOptions *workflow.ActivityOptions
	}

	WorkflowExecutionData struct {
		ID                     uint64
		WorkflowID             string
		CurrState              string
		RunID                  string
		SignalEvent            string
		activityExecutionInfos []ActivityExecutionInfo          // sequential pipeline
		branchActivities       map[string]ActivityExecutionInfo // branch-only targets (compensation, etc.)
		StartedAt              time.Time
		CompletedAt            time.Time

		temporalClient Temporal
	}

	WorkflowExecution interface {
		// Execute runs the sequential activity pipeline, threading state through each activity.
		Execute(ctx workflow.Context, state interface{}) error

		// AddTransitionActivity registers an activity with the Temporal worker and adds it
		// to the sequential execution pipeline. Activities run in the order they are added.
		AddTransitionActivity(activityName string, signalName string, activityFn interface{})

		// AddTransitionActivityWithOptions registers an activity with the Temporal worker and adds it
		// to the sequential execution pipeline. Activities run in the order they are added.
		// It is used to add an activity with options to the sequential execution pipeline.
		AddTransitionActivityWithOptions(activityName string, signalName string, activityFn interface{}, options *workflow.ActivityOptions)

		// AddBranchActivity registers an activity that is only reachable via branching.
		// It is NOT part of the sequential pipeline — it only runs when another activity
		// sets NextActivity to this activity's name.
		// Use this for compensation, rollback, or alternative execution paths.
		AddBranchActivity(activityName string, activityFn interface{})

		// RegisterWorkflow registers a workflow with the Temporal worker.
		RegisterWorkflow(name string, fn interface{})

		// GetWorkflowExecutionData gets the workflow execution data.
		// It is used to get the workflow execution data from the Temporal server.
		GetWorkflowExecutionData(wfCtx workflow.Context, runID string, result interface{}) error

		// StartWorkflow starts a new workflow execution and returns the run ID.
		// It is used to start a new workflow execution and returns the run ID.
		StartWorkflow(ctx context.Context, opts StartWorkflowOptions, workflowFn interface{}, args ...interface{}) (client.WorkflowRun, error)

		// GetWorkflowResult gets the workflow result from the Temporal server.
		// It is used to get the workflow result from the Temporal server.
		GetWorkflowResult(ctx context.Context, workflowID string, runID string, result interface{}) error
	}
)

// Execute runs the sequential activity pipeline, threading state through each activity.
// If the state implements Navigable and an activity sets NextActivity, Execute branches
// to that activity (which must be registered via AddBranchActivity). After the branch
// chain completes, Execute returns — it does NOT resume the sequential pipeline.
func (w *WorkflowExecutionData) Execute(ctx workflow.Context, state interface{}) error {
	navigable, hasBranching := state.(Navigable)

	for i := 0; i < len(w.activityExecutionInfos); i++ {
		info := w.activityExecutionInfos[i]

		if err := w.runActivity(ctx, info, state); err != nil {
			return err
		}

		// Check if the activity wants to branch
		if hasBranching {
			if next := navigable.GetNextActivity(); next != "" {
				navigable.SetNextActivity("")
				return w.executeBranch(ctx, next, state, navigable)
			}
		}
	}

	return nil
}

func (w *WorkflowExecutionData) StartWorkflow(ctx context.Context, opts StartWorkflowOptions, workflowFn interface{}, args ...interface{}) (client.WorkflowRun, error) {
	return w.temporalClient.StartWorkflow(ctx, opts, workflowFn, args...)
}

func (w *WorkflowExecutionData) GetWorkflowResult(ctx context.Context, workflowID string, runID string, result interface{}) error {
	return w.temporalClient.GetWorkflowResult(ctx, workflowID, runID, result)
}

// executeBranch follows a chain of branch activities. Each branch activity
// can set NextActivity to chain to another branch activity. When an activity
// does not set NextActivity, the chain ends and executeBranch returns.
func (w *WorkflowExecutionData) executeBranch(ctx workflow.Context, activityName string, state interface{}, navigable Navigable) error {
	current := activityName
	for current != "" {
		info, ok := w.branchActivities[current]
		if !ok {
			return fmt.Errorf("branch activity %q not found", current)
		}

		if err := w.runActivity(ctx, info, state); err != nil {
			return err
		}

		current = navigable.GetNextActivity()
		navigable.SetNextActivity("")
	}
	return nil
}

// runActivity executes a single activity and handles its signal if present.
func (w *WorkflowExecutionData) runActivity(ctx workflow.Context, info ActivityExecutionInfo, state interface{}) error {
	activityCtx := ctx

	if info.ActivityOptions != nil {
		activityCtx = workflow.WithActivityOptions(ctx, *info.ActivityOptions)
	}

	future := workflow.ExecuteActivity(activityCtx, info.ActivityName, state)
	if err := future.Get(ctx, state); err != nil {
		return fmt.Errorf("activity %s failed: %w", info.ActivityName, err)
	}

	if info.SignalName != "" {
		if err := ExecuteChildWorkflow(ctx, info.SignalName, state); err != nil {
			return fmt.Errorf("child workflow for activity %s failed: %w", info.ActivityName, err)
		}
	}
	return nil
}

// AddTransitionActivity registers an activity with the Temporal worker and adds it
// to the sequential execution pipeline. Activities run in the order they are added.
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

// AddBranchActivity registers an activity that is only reachable via branching.
// It is NOT part of the sequential pipeline — it only runs when another activity
// sets NextActivity to this activity's name.
// Use this for compensation, rollback, or alternative execution paths.
func (w *WorkflowExecutionData) AddBranchActivity(activityName string, activityFn interface{}) {
	w.temporalClient.RegisterActivity(ActivityDefinition{
		Name: activityName,
		Fn:   activityFn,
	})

	w.branchActivities[activityName] = ActivityExecutionInfo{
		ActivityName: activityName,
	}
}

// RegisterWorkflow registers a workflow with the Temporal worker.
func (w *WorkflowExecutionData) RegisterWorkflow(name string, fn interface{}) {
	w.temporalClient.RegisterWorkflow(WorkflowDefinition{
		Name: name,
		Fn:   fn,
	})
}

// GetWorkflowExecutionData gets the workflow execution data.
// It is used to get the workflow execution data from the Temporal server.
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

func (w *WorkflowExecutionData) AddTransitionActivityWithOptions(activityName string, signalName string, activityFn interface{}, options *workflow.ActivityOptions) {
	w.temporalClient.RegisterActivity(ActivityDefinition{
		Name: activityName,
		Fn:   activityFn,
	})

	w.activityExecutionInfos = append(w.activityExecutionInfos, ActivityExecutionInfo{
		ActivityName:    activityName,
		SignalName:      signalName,
		ActivityOptions: options,
	})
}

// NewWorkflowExecution creates a new WorkflowExecution.
// It is used to create a new WorkflowExecution.
func NewWorkflowExecution(
	temporalClient Temporal,
) WorkflowExecution {
	return &WorkflowExecutionData{
		activityExecutionInfos: make([]ActivityExecutionInfo, 0),
		branchActivities:       make(map[string]ActivityExecutionInfo),
		temporalClient:         temporalClient,
	}
}
