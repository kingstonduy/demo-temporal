package workflow

import (
	"time"

	activity "demo-temporal/activity"

	"go.temporal.io/sdk/workflow"
)

func AsyncWorkFlow(ctx workflow.Context) error {

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

	future := workflow.ExecuteActivity(ctx, activity.GetOcbInfo)

	var flag bool
	err := workflow.ExecuteActivity(ctx, activity.Withdraw, nil).Get(ctx, &flag)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, activity.UserInputOtp, flag).Get(ctx, &flag)

	for flag == false {
		err = workflow.ExecuteActivity(ctx, activity.ResendOtp, nil).Get(ctx, nil)
		if err != nil {
			return err
		}

		err = workflow.ExecuteActivity(ctx, activity.UserInputOtp, flag).Get(ctx, &flag)
		if err != nil {
			return err
		}
	}

	if flag == true {
		err = workflow.ExecuteActivity(ctx, activity.Notification, flag).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	err = future.Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
