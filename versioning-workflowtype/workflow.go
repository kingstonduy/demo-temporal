package versioning_workflowtype

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

var timeout time.Duration = 30

func SimpleWorkflow(ctx workflow.Context) error {

	// retryPolicy := &temporal.RetryPolicy{
	// 	InitialInterval: time.Second,
	// 	MaximumAttempts: 1, // unlimited retries
	// }

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * timeout,
		HeartbeatTimeout:    time.Second * timeout,
		// RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)
	workflow.ExecuteActivity(ctx, GetInformation).Get(ctx, nil)
	return nil
}

func SimpleWorkflow1(ctx workflow.Context) error {

	// retryPolicy := &temporal.RetryPolicy{
	// 	InitialInterval: time.Second,
	// 	MaximumAttempts: 1, // unlimited retries
	// }

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * timeout,
		HeartbeatTimeout:    time.Second * timeout,
		// RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)
	workflow.ExecuteActivity(ctx, GetInformation1).Get(ctx, nil)
	return nil
}
