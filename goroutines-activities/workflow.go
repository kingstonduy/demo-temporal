package goroutines_activities

import (
	"time"

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
	err := workflow.ExecuteActivity(ctx, Withdraw, nil).Get(ctx, &flag)
	if err != nil {
		return err
	}

	wg := workflow.NewWaitGroup(ctx)
	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err = workflow.ExecuteActivity(ctx, UserInputOtp, flag).Get(ctx, &flag)
		if err != nil {
			return
		}

		if flag == true {
			err = workflow.ExecuteActivity(ctx, Notification, flag).Get(ctx, nil)
			if err != nil {
				return
			}
		}

		defer wg.Done()
	})

	err = workflow.ExecuteActivity(ctx, LongAcitivity, nil).Get(ctx, nil)
	if err != nil {
		return err
	}

	wg.Wait(ctx)
	return nil
}
