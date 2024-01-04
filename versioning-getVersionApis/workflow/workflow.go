package workflow

import (
	"demo-temporal/activity"
	"time"

	"go.temporal.io/sdk/workflow"
)

func SimpleWorkflow(ctx workflow.Context) error {

	// retryPolicy := &temporal.RetryPolicy{
	// 	InitialInterval: time.Second,
	// 	MaximumAttempts: 1, // unlimited retries
	// }

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 15,
		HeartbeatTimeout:    time.Second * 15,
		// RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// workflow.ExecuteActivity(ctx, activity.GetInformation).Get(ctx, nil)

	version := workflow.GetVersion(ctx, "UpdateInformation1", workflow.DefaultVersion, 1)

	if version == workflow.DefaultVersion {
		// in ra default
		workflow.ExecuteActivity(ctx, activity.GetInformation).Get(ctx, nil)
	}

	if version == 1 {
		// in ra version 1
		workflow.ExecuteActivity(ctx, activity.GetInformation1).Get(ctx, nil)
	}

	// if version == 2 {
	// 	workflow.ExecuteActivity(ctx, activity.GetInformation2).Get(ctx, nil)
	// }

	return nil
}
