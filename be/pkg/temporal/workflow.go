package temporal_client

import (
	"context"
	"fmt"
	"log"
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
	NavigatableActivity interface {
		GetCurrentActivity() string
		SetCurrentActivity(name string)
		GetNextActivity() string
		SetNextActivity(name string)
	}

	SignalActivity struct {
	}

	ActivityExecutionInfo struct {
		ActivityName    string
		SignalName      string
		ActivityFn      interface{}
		ActivityOptions *workflow.ActivityOptions
		NextActivities  []string
	}

	WorkflowExecutionData struct {
		ID                     uint64
		WorkflowID             string
		CurrState              string
		RunID                  string
		SignalEvent            string
		activityExecutionInfos []ActivityExecutionInfo
		branchActivities       map[string]ActivityExecutionInfo // branch-only targets (compensation, etc.)
		activity               map[string]ActivityExecutionInfo
		firstActivity          string
		StartedAt              time.Time
		CompletedAt            time.Time

		temporalClient Temporal
	}

	WorkflowExecution interface {
		// Execute runs the sequential activity pipeline, threading state through each activity.
		Execute(ctx workflow.Context, executionData interface{}) error

		// AddTransitionActivity registers an activity with the Temporal worker and adds it
		// to the sequential execution pipeline. Activities run in the order they are added.
		AddTransitionActivity(activityName string, signalName string, activityFn interface{})

		// AddTransitionActivityWithOptions registers an activity with the Temporal worker and adds it
		// to the sequential execution pipeline. Activities run in the order they are added.
		// It is used to add an activity with options to the sequential execution pipeline.
		AddTransitionActivityWithOptions(activityName string, signalName string, activityFn interface{}, options *workflow.ActivityOptions, nextActivities ...string)

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

		// StartChildWorkflow starts a new child workflow execution and returns the run ID.
		// It is used to start a new child workflow execution and returns the run ID.
		StartChildWorkflow(ctx workflow.Context, workflowID string, signalName string, request interface{}, result interface{}) error

		//SignalWorkflow signals a workflow.
		SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error

		//Goroutine workflow run
		Goroutine(ctx workflow.Context, goroutineFn func(ctx workflow.Context))
	}
)

// Execute runs the sequential activity pipeline, threading state through each activity.
// If the state implements Navigable and an activity sets NextActivity, Execute branches
// to that activity (which must be registered via AddBranchActivity). After the branch
// chain completes, Execute returns — it does NOT resume the sequential pipeline.
func (w *WorkflowExecutionData) Execute(ctx workflow.Context, executionData interface{}) error {
	navigable, _ := executionData.(NavigatableActivity)
	currActivity := w.activity[w.firstActivity]

	for currActivity.NextActivities != nil {
		if err := w.runActivity(ctx, currActivity, executionData, navigable); err != nil {
			return err
		}

		nextActivity, err := w.getNextActivity(currActivity, navigable.GetNextActivity())
		if err != nil {
			return err
		}

		currActivity = *nextActivity
	}

	return nil
}

func (w *WorkflowExecutionData) getNextActivity(currActivity ActivityExecutionInfo, nextActivity string) (*ActivityExecutionInfo, error) {
	log.Println("next activity", nextActivity)
	for _, info := range currActivity.NextActivities {
		if info == nextActivity {
			result := w.activity[info]
			return &result, nil
		}
	}
	return nil, fmt.Errorf("next activity not found")
}

// StartWorkflow starts a new workflow execution and returns the run ID.
func (w *WorkflowExecutionData) StartWorkflow(ctx context.Context, opts StartWorkflowOptions, workflowFn interface{}, args ...interface{}) (client.WorkflowRun, error) {
	return w.temporalClient.StartWorkflow(ctx, opts, workflowFn, args...)
}

// GetWorkflowResult gets the workflow result from the Temporal server.
func (w *WorkflowExecutionData) GetWorkflowResult(ctx context.Context, workflowID string, runID string, result interface{}) error {
	return w.temporalClient.GetWorkflowResult(ctx, workflowID, runID, result)
}

func (w *WorkflowExecutionData) SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error {
	return w.temporalClient.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
}

// executeBranch follows a chain of branch activities. Each branch activity
// can set NextActivity to chain to another branch activity. When an activity
// does not set NextActivity, the chain ends and executeBranch returns.
func (w *WorkflowExecutionData) executeBranch(ctx workflow.Context, activityName string, executionData interface{}, navigable NavigatableActivity) error {
	current := activityName
	for current != "" {
		info, ok := w.branchActivities[current]
		if !ok {
			return fmt.Errorf("branch activity %q not found", current)
		}

		if err := w.runActivity(ctx, info, executionData, navigable); err != nil {
			return err
		}

		current = navigable.GetCurrentActivity()
		navigable.SetCurrentActivity("")
	}
	return nil
}

// runActivity executes a single activity and handles its signal if present.
func (w *WorkflowExecutionData) runActivity(ctx workflow.Context, info ActivityExecutionInfo, executionData interface{}, navigable NavigatableActivity) error {
	activityCtx := ctx

	if info.ActivityOptions != nil {
		activityCtx = workflow.WithActivityOptions(ctx, *info.ActivityOptions)
	}

	future := workflow.ExecuteActivity(activityCtx, info.ActivityFn, executionData)
	if err := future.Get(ctx, executionData); err != nil {
		return fmt.Errorf("activity %s failed: %w", info.ActivityName, err)
	}

	if navigable.GetCurrentActivity() != "" {
		for _, nextActivity := range info.NextActivities {
			if navigable.GetNextActivity() == nextActivity {
				navigable.SetCurrentActivity(nextActivity)
				return nil
			}
		}
	}

	if info.SignalName != "" {
		if err := w.StartChildWorkflow(ctx, w.WorkflowID, info.SignalName, executionData, executionData); err != nil {
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

func (w *WorkflowExecutionData) AddTransitionActivityWithOptions(activityName string, signalName string, activityFn interface{}, options *workflow.ActivityOptions, nextActivities ...string) {
	w.temporalClient.RegisterActivity(ActivityDefinition{
		Name: activityName,
		Fn:   activityFn,
	})

	//initiate first activity
	if len(w.activity) == 0 {
		w.firstActivity = activityName
	}

	w.activity[activityName] = ActivityExecutionInfo{
		ActivityName:    activityName,
		SignalName:      signalName,
		ActivityFn:      activityFn,
		ActivityOptions: options,
		NextActivities:  nextActivities,
	}

	for _, nextActivity := range nextActivities {
		w.activity[nextActivity] = ActivityExecutionInfo{}
	}
}

func (w *WorkflowExecutionData) StartChildWorkflow(ctx workflow.Context, workflowID string, signalName string, request interface{}, result interface{}) error {
	childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: workflowID,
	})

	childWorkflowRun := workflow.ExecuteChildWorkflow(childCtx, signalName)
	var workflowExecution workflow.Execution
	if err := childWorkflowRun.GetChildWorkflowExecution().Get(ctx, &workflowExecution); err != nil {
		return fmt.Errorf("failed to get child workflow execution: %w", err)
	}

	sigFuture := workflow.SignalExternalWorkflow(ctx, workflowExecution.ID, workflowExecution.RunID, signalName, request)
	if err := sigFuture.Get(ctx, nil); err != nil {
		return fmt.Errorf("failed to signal child workflow: %w", err)
	}

	if err := childWorkflowRun.Get(childCtx, result); err != nil {
		return fmt.Errorf("failed to get child workflow result: %w", err)
	}

	return nil
}

func (w *WorkflowExecutionData) Goroutine(ctx workflow.Context, goroutineFn func(ctx workflow.Context)) {
	workflow.Go(ctx, goroutineFn)
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
		activity:               make(map[string]ActivityExecutionInfo),
	}
}
