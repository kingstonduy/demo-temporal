package workflow

import (
	"time"

	activity "demo-temporal/activity"

	"go.temporal.io/sdk/workflow"
)

func AsyncWorkFlow1(ctx workflow.Context) error {

	// retryPolicy := &temporal.RetryPolicy{
	// 	InitialInterval: time.Second,
	// 	MaximumAttempts: 1, // unlimited retries
	// }

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 1,
		HeartbeatTimeout:    time.Minute * 1,
		// RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var flag bool
	err := workflow.ExecuteActivity(ctx, activity.Withdraw, nil).Get(ctx, &flag)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, activity.UserInputOtp, flag).Get(ctx, &flag)
	if err != nil {
		return err
	}

	// getting blocked activity
	future := workflow.ExecuteActivity(ctx, activity.LongAcitivity, nil)

	if flag == true {
		err = workflow.ExecuteActivity(ctx, activity.Notification, flag).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	future.Get(ctx, nil)

	return nil
}
