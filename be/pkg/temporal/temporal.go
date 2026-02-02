package temporal_client

import (
	"context"
	"fmt"

	"github.com/RandySteven/CafeConnect/be/configs"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type (
	Temporal interface {
		RegisterWorkflow(workflow interface{})
		RegisterActivity(activity interface{}, activityName string)
		StartWorkflow(ctx context.Context, workflowID string, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
	}

	temporalClient struct {
		worker worker.Worker
		client client.Client
	}
)

func (t *temporalClient) RegisterActivity(activityFn interface{}, activityName string) {
	t.worker.RegisterActivityWithOptions(activityFn, activity.RegisterOptions{
		Name: activityName,
	})
}

func (t *temporalClient) RegisterWorkflow(workflow interface{}) {
	t.worker.RegisterWorkflow(workflow)
}

func (t *temporalClient) StartWorkflow(ctx context.Context, workflowID string, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	return t.client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID: workflowID,
	}, workflow, args...)
}

func NewTemporalClient(config *configs.Config) (*temporalClient, error) {
	client, err := client.NewClient(client.Options{
		HostPort:  fmt.Sprintf("%s:%s", config.Config.Temporal.Host, config.Config.Temporal.Port),
		Namespace: config.Config.Temporal.Namespace,
	})
	if err != nil {
		return nil, err
	}

	var workerOptions = worker.Options{}
	if config.Config.Temporal.WorkerOptions != nil {
		workerOptions = worker.Options{
			MaxConcurrentActivityExecutionSize:      config.Config.Temporal.WorkerOptions.MaxConcurrentActivityExecutionSize,
			WorkerActivitiesPerSecond:               config.Config.Temporal.WorkerOptions.WorkerActivitiesPerSecond,
			MaxConcurrentLocalActivityExecutionSize: config.Config.Temporal.WorkerOptions.MaxConcurrentLocalActivityExecutionSize,
			WorkerLocalActivitiesPerSecond:          config.Config.Temporal.WorkerOptions.WorkerLocalActivitiesPerSecond,
		}
	}

	return &temporalClient{
		client: client,
		worker: worker.New(client, config.Config.Temporal.TaskQueue, workerOptions),
	}, nil
}
