package workflow

import (
	"demo-temporal/activity"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

func SimpleWorkflowV1(ctx workflow.Context) error {
	fmt.Println("HIHI")
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 15,
		HeartbeatTimeout:    time.Second * 15,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, activity.GetInformation).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func SimpleWorkflowV2(ctx workflow.Context) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 15,
		HeartbeatTimeout:    time.Second * 15,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, activity.GetInformation1).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func SimpleWorkflowV3(ctx workflow.Context) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 15,
		HeartbeatTimeout:    time.Second * 15,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, activity.GetInformation2).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
